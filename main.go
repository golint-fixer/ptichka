package main

import (
	"fmt"
	"log"
	"os/user"
	"sort"
	"time"
)

// Tweet is a simplified anaconda.Tweet.
type Tweet struct {
	IDStr          string
	UserScreenName string
	Date           time.Time
	Text           string
	Medias         []string
}

// TweetsByDate is a slice of Tweet
// with ability to be sorted by date from older to newer.
type TweetsByDate []Tweet

// <https://github.com/wskinner/anaconda/commit/d0c12d8fba671d7d5ce27d3abd1809aedcc59195>,
// <http://nerdyworm.com/blog/2013/05/15/sorting-a-slice-of-structs-in-go/>.

// Len is the number of elements in the collection.
func (a TweetsByDate) Len() int { return len(a) }

// Less reports whether the element with index i should sort before
// the element with index j.
func (a TweetsByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Swap swaps the elements with indexes i and j.
func (a TweetsByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func main() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	config, err := loadConfig(fmt.Sprintf("%s/.tgtmrc.toml", currentUser.HomeDir))
	if err != nil {
		log.Fatal(err)
	}

	oldIds, err := loadCache(config.CacheFile)
	if err != nil {
		log.Fatal(err)
	}

	anacondaTweets, err := fetchTweets(config)
	if err != nil {
		log.Fatal(err)
	}

	tweets := anacondaTweets.toTweets()
	sort.Sort(tweets)

	newIds := filterTweets(tweets, oldIds, config)

	if err := saveCache(config.CacheFile, newIds); err != nil {
		log.Fatal(err)
	}
}
