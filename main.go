package main

import (
	"bytes"
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
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

			subject := config.Label +
				"@" + currentTweet.User.ScreenName +
				" " +
				createdAt

			print(subject)
			print("\n")
			print("\n")

			body, err := tweetBody(tweet{
				currentTweet.IdStr,
				currentTweet.User.ScreenName,
				currentTweet.Text})
			if err != nil {
				panic(err)
			}
			print(body)
			print("\n")

			print("\n")
			print("\n")
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
