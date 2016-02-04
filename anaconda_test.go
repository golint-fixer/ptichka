package main

import (
	"encoding/json"
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
        {"Media_url_https": ":this://is.wrong/url/should/be/rejected"},
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

	var anacondaTweets AnacondaTweets
	json.Unmarshal(jsonBlob, &anacondaTweets)
	tweets := anacondaTweets.toTweets()

	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	referenceTweets := TweetsByDate{
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

	for i := range referenceTweets {
		referenceTweet := referenceTweets[i]
		tweet := tweets[i]

		if referenceTweet.ID != tweet.ID {
			t.Errorf(
				"ReferenceTweet.ID{%q} =! anaconda.Tweet.IdStr{%q}",
				referenceTweet.ID,
				tweet.ID)
		}

		if !referenceTweet.Date.Equal(tweet.Date) {
			t.Errorf(
				"ReferenceTweet.Date{%q} =! anaconda.Tweet.createdAt{%q}",
				referenceTweet.Date,
				tweet.Date)
		}

		if referenceTweet.User != tweet.User {
			t.Errorf(
				"ReferenceTweet.User{%q} =! anaconda.Tweet.User.ScreenName{%q}",
				referenceTweet.User,
				tweet.User)
		}

		if referenceTweet.Text != tweet.Text {
			t.Errorf(
				"ReferenceTweet.Text{%q} =! anaconda.Tweet.Text{%q}",
				referenceTweet.Text,
				tweet.Text)
		}

		for j := range referenceTweet.Medias {
			referenceMedia := referenceTweet.Medias[j]
			media := tweet.Medias[j]

			if referenceMedia != media {
				t.Errorf(
					"ReferenceTweet.Medias[%d]{%q} =! anaconda.Tweet.ExtendedEntities.Media[%d]{%q}",
					j,
					referenceMedia,
					j,
					media)
			}
		}
	}
}
