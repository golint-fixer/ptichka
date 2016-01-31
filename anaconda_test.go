package main

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"log"
	"testing"
	"time"
)

func TestToTweets(t *testing.T) {
	var jsonBlob = []byte(`[
  {
    "id_str": "111111111111111111",
    "created_at": "Thu Jan 01 01:00:00 +0000 1970",
    "user": {"screen_name": "johndoe"},
    "text": "RT @ivanivanov: Hello, World!"
  },
  {
    "id_str": "222222222222222222",
    "created_at": "Thu Jan 02 02:00:00 +0000 1970",
    "user": {"screen_name": "johndoe"},
    "text": "foo bar",
    "extended_entities": {
      "Media": [
        {"Media_url_https": "https://pbs.twimg.com/media/qwertyuiopasdfg.jpg"},
        {"Media_url_https": "https://pbs.twimg.com/media/hjklzxcvbnmqwer.png"}
      ]
    }
  },
  {
    "id_str": "333333333333333333",
    "created_at": "Thu Jan 03 03:00:00 +0000 1970",
    "user": {"screen_name": "gunterschmidt"},
    "text": "baz xyz",
    "extended_entities": {
      "Media": [
        {"Media_url_https": "https://pbs.twimg.com/media/tyuiopasdfghjkl.gif"}
      ]
    }
  }
]`)

	var anacondaTweets []anaconda.Tweet
	json.Unmarshal(jsonBlob, &anacondaTweets)

	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	tweets := TweetsByDate{
		Tweet{
			ID:   "111111111111111111",
			Date: time.Date(1970, 1, 1, 1, 0, 0, 0, utc),
			User: "johndoe",
			Text: "RT @ivanivanov: Hello, World!"},
		Tweet{
			ID:   "222222222222222222",
			Date: time.Date(1970, 1, 2, 2, 0, 0, 0, utc),
			User: "johndoe",
			Text: "foo bar",
			Medias: []string{
				"https://pbs.twimg.com/media/qwertyuiopasdfg.jpg",
				"https://pbs.twimg.com/media/hjklzxcvbnmqwer.png"}},
		Tweet{
			ID:   "333333333333333333",
			Date: time.Date(1970, 1, 3, 3, 0, 0, 0, utc),
			User: "gunterschmidt",
			Text: "baz xyz",
			Medias: []string{
				"https://pbs.twimg.com/media/tyuiopasdfghjkl.gif"}}}

	for i := range tweets {
		tweet := tweets[i]
		anacondaTweet := anacondaTweets[i]

		if tweet.ID != anacondaTweet.IdStr {
			t.Errorf(
				"Tweet.ID{%q} =! anaconda.Tweet.IdStr{%q}",
				tweet.ID,
				anacondaTweet.IdStr)
		}

		createdAt, err := time.Parse(time.RubyDate, anacondaTweet.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		if !tweet.Date.Equal(createdAt) {
			t.Errorf(
				"Tweet.Date{%q} =! anaconda.Tweet.createdAt{%q}",
				tweet.Date,
				createdAt)
		}

		if tweet.User != anacondaTweet.User.ScreenName {
			t.Errorf(
				"Tweet.User{%q} =! anaconda.Tweet.User.ScreenName{%q}",
				tweet.User,
				anacondaTweet.User.ScreenName)
		}

		if tweet.Text != anacondaTweet.Text {
			t.Errorf(
				"Tweet.Text{%q} =! anaconda.Tweet.Text{%q}",
				tweet.Text,
				anacondaTweet.Text)
		}

		for j := range tweet.Medias {
			media := tweet.Medias[j]
			anacondaMedia := anacondaTweet.ExtendedEntities.Media[j].Media_url_https

			if media != anacondaMedia {
				t.Errorf(
					"Tweet.Medias[%d]{%q} =! anaconda.Tweet.ExtendedEntities.Media[%d]{%q}",
					j,
					media,
					j,
					anacondaMedia)
			}
		}
	}
}
