package main

import (
	// "github.com/ChimeraCoder/anaconda"
	"log"
)

func main() {
	config, err := loadConfig(".twitmarc.toml")
	if err != nil {
		log.Fatal(err)
	}

	ids, err := loadCache(config.CacheFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(ids) > 0 {
		print(ids[0])
	} else {
		print("Have no ids(")
	}

	// tweets, err := fetchTweets(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// j, err := json.Marshal(tweets)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ioutil.WriteFile("foobar.json", j, 0644)
}

// func fetchTweets(config *config) ([]anaconda.Tweet, error) {
// 	anaconda.SetConsumerKey(config.Twitter.ConsumerKey)
// 	anaconda.SetConsumerSecret(config.Twitter.ConsumerSecret)
// 	api := anaconda.NewTwitterApi(config.Twitter.AccessToken, config.Twitter.AccessTokenSecret)
// 	tweets, err := api.GetHomeTimeline(nil)
// 	return tweets, err
// }
