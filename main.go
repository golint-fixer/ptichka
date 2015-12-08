package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
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
				currentTweet.IdStr,
				currentTweet.User.ScreenName,
				currentTweet.Text})
			if err != nil {
				panic(err)
			}

			from := mail.Address{"", config.Mail.From.Address}
			to := mail.Address{config.Mail.To.Name, config.Mail.To.Address}

			// Setup headers
			headers := make(map[string]string)
			headers["From"] = from.String()
			headers["To"] = to.String()
			headers["Subject"] = subject

			// Setup message
			message := ""
			for k, v := range headers {
				message += fmt.Sprintf("%s: %s\r\n", k, v)
			}
			message += "\r\n" + body

			// Connect to the SMTP Server
			servername := fmt.Sprintf(
				"%s:%d",
				config.Mail.SMTP.Address,
				config.Mail.SMTP.Port)

			host, _, _ := net.SplitHostPort(servername)

			auth := smtp.PlainAuth(
				"",
				config.Mail.SMTP.UserName,
				config.Mail.SMTP.Password,
				host)

			// TLS config
			tlsconfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host,
			}

			// Here is the key, you need to call tls.Dial instead of smtp.Dial
			// for smtp servers running on 465 that require an ssl connection
			// from the very beginning (no starttls)
			conn, err := tls.Dial("tcp", servername, tlsconfig)
			if err != nil {
				log.Panic(err)
			}
			defer conn.Close()

			c, err := smtp.NewClient(conn, host)
			if err != nil {
				log.Panic(err)
			}
			defer c.Quit()

			// Auth
			if err = c.Auth(auth); err != nil {
				log.Panic(err)
			}

			// To && From
			if err = c.Mail(from.Address); err != nil {
				log.Panic(err)
			}

			if err = c.Rcpt(to.Address); err != nil {
				log.Panic(err)
			}

			// Data
			w, err := c.Data()
			if err != nil {
				log.Panic(err)
			}
			defer w.Close()

			_, err = w.Write([]byte(message))
			if err != nil {
				log.Panic(err)
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
