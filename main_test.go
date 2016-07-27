package main

import (
	"log"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	got := TweetsByDate{
		Tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)}}

	wont := TweetsByDate{
		Tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)}}

	sort.Sort(got)

	if !reflect.DeepEqual(got, wont) {
		t.Errorf("%v != %v", got, wont)
	}
}

func TestTweetBody(t *testing.T) {
	got, err := tweetBody(Tweet{
		IDStr:          "1234",
		UserScreenName: "johndoe",
		Text:           "Hello &amp; world!"})

	if err != nil {
		t.Errorf("Error on tweetBody(tweet{...}): %v", err)
	}

	wont := `@johndoe

Hello & world!

https://twitter.com/johndoe/status/1234`

	if got != wont {
		t.Errorf(`tweetBody(tweet{IDStr: "1234", UserScreenName: "johndoe", Text: "Hello &amp; world!"})
get:
%v

wont:
%v`, got, wont)
	}
}
