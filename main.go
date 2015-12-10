package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"net/mail"
	"text/template"
	"time"
)

func main() {
	config, err := loadConfig(".twitmarc.toml")
	if err != nil {
		log.Fatal(err)
	}

	oldIds, err := loadCache(config.CacheFile)
	if err != nil {
		log.Fatal(err)
	}

	tweets, err := fetchTweets(config)
	if err != nil {
		log.Fatal(err)
	}
	newIds := oldIds
	for _, currentTweet := range tweets {
		if !contains(oldIds, currentTweet.IdStr) {
			newIds = append(newIds, currentTweet.IdStr)

			t, err := time.Parse(time.RubyDate, currentTweet.CreatedAt)
			if err != nil {
				log.Fatal(err)
			}
			createdAt := t.Format("2006-01-02 15:04 -0700")

			subject := fmt.Sprintf(
				"%s@%s %s",
				config.Label,
				currentTweet.User.ScreenName,
				createdAt)

			if config.Verbose {
				print("Sending: " + subject)
			}

			body, err := tweetBody(tweet{
				ID:   currentTweet.IdStr,
				User: currentTweet.User.ScreenName,
				Text: currentTweet.Text})
			if err != nil {
				panic(err)
			}

			from := mail.Address{
				Name:    config.Mail.From.Name,
				Address: config.Mail.From.Address}
			to := mail.Address{
				Name:    config.Mail.To.Name,
				Address: config.Mail.To.Address}

			message := gomail.NewMessage()
			message.SetHeader("From", from.String())
			message.SetHeader("To", to.String())
			message.SetHeader("Subject", subject)
			message.SetBody("text/plain", body)
			// message.Attach("/home/Alex/lolcat.jpg")

			dialer := gomail.NewPlainDialer(
				config.Mail.SMTP.Address,
				config.Mail.SMTP.Port,
				config.Mail.SMTP.UserName,
				config.Mail.SMTP.Password)

			if err := dialer.DialAndSend(message); err != nil {
				panic(err)
			}

			if config.Verbose {
				print("\n")
			}
		}
	}
	if err := saveCache("xyz.json", newIds); err != nil {
		log.Fatal(err)
	}
}

type tweet struct {
	ID   string
	User string
	Text string
}

func tweetBody(tweet tweet) (string, error) {
	tmpl, err := template.New("tweet").Parse(
		`@{{.User}}

{{.Text}}

https://twitter.com/{{.User}}/status/{{.ID}}`)
	if err != nil {
		panic(err)
	}

	var x bytes.Buffer

	err = tmpl.Execute(&x, tweet)

	return x.String(), err
}

func contains(ids []string, id string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

func fetchTweets(config *config) ([]anaconda.Tweet, error) {
	// anaconda.SetConsumerKey(config.Twitter.ConsumerKey)
	// anaconda.SetConsumerSecret(config.Twitter.ConsumerSecret)
	// api := anaconda.NewTwitterApi(
	// 	config.Twitter.AccessToken,
	// 	config.Twitter.AccessTokenSecret)
	// tweets, err := api.GetHomeTimeline(nil)
	// return tweets, err

	foobar, err := ioutil.ReadFile("foobar.json")
	if err != nil {
		log.Fatal(err)
	}

	var tweets []anaconda.Tweet
	err = json.Unmarshal(foobar, &tweets)
	return tweets, nil
}
