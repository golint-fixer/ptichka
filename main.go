package main

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
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
			print("aaaaaaaaaaaaa: ")
		}
		print(tweet.IdStr)
		print("\n")
	}
	if err := saveCache("xyz.json", newIds); err != nil {
		log.Fatal(err)
	}
	if len(newIds) > len(oldIds) {
		print(len(newIds))
		print("\n")
		print(len(oldIds))
		print("\n")
		print("bbbbbbbbbbbb: ")
		print(newIds[0])
		print("\n")
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
