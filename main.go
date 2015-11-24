package main

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
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
	for _, tweet := range tweets {
		if !contains(oldIds, tweet.IdStr) {
			newIds = append(newIds, tweet.IdStr)

			user := "@" + tweet.User.ScreenName

			t, err := time.Parse(time.RubyDate, tweet.CreatedAt)
			if err != nil {
				log.Fatal(err)
			}
			createdAt := t.Format("2006-01-02 15:04 -0700")

			print("id: ")
			print(tweet.IdStr)
			print("\n")
			print("\n")

			subject := config.Label + user + " " + createdAt
			print(subject)
			print("\n")
			print("\n")

			print(user)
			print("\n")
			print("\n")

			print(tweet.Text)
			print("\n")
			print("\n")

			body := "https://twitter.com/" +
				tweet.User.ScreenName +
				"/status/" +
				tweet.IdStr
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
