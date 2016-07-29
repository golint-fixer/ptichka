package ptichka

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
    "text": "RT @ivanivanov: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqu...",
    "retweeted_status": {
      "user": {
        "screen_name": "ivanivanov"
      },
      "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim"
    }
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

	var anacondaTweets anacondaTweets
	err := json.Unmarshal(jsonBlob, &anacondaTweets)
	if err != nil {
		log.Fatal(err)
	}
	tweets, err := anacondaTweets.toTweets()
	if err != nil {
		log.Fatal(err)
	}

	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	referenceTweets := TweetsByDate{
		Tweet{
			IDStr:          "111111111111111111",
			Date:           time.Date(1970, 1, 1, 1, 0, 0, 0, utc),
			UserScreenName: "johndoe",
			Text:           "RT @ivanivanov: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim"},
		Tweet{
			IDStr:          "222222222222222222",
			Date:           time.Date(1970, 1, 2, 2, 0, 0, 0, utc),
			UserScreenName: "johndoe",
			Text:           "foo bar",
			Medias: []string{
				"https://pbs.twimg.com/media/qwertyuiopasdfg.jpg",
				"https://pbs.twimg.com/media/hjklzxcvbnmqwer.png"}},

		Tweet{
			IDStr:          "333333333333333333",
			Date:           time.Date(1970, 1, 3, 3, 0, 0, 0, utc),
			UserScreenName: "gunterschmidt",
			Text:           "baz xyz",
			Medias: []string{
				"https://pbs.twimg.com/media/tyuiopasdfghjkl.gif"}}}

	for i := range referenceTweets {
		referenceTweet := referenceTweets[i]
		tweet := tweets[i]

		if referenceTweet.IDStr != tweet.IDStr {
			t.Errorf(
				"ReferenceTweet.IDStr{%q} =! anaconda.Tweet.IdStr{%q}",
				referenceTweet.IDStr,
				tweet.IDStr)
		}

		if !referenceTweet.Date.Equal(tweet.Date) {
			t.Errorf(
				"ReferenceTweet.Date{%q} =! anaconda.Tweet.createdAt{%q}",
				referenceTweet.Date,
				tweet.Date)
		}

		if referenceTweet.UserScreenName != tweet.UserScreenName {
			t.Errorf(
				"ReferenceTweet.UserScreenName{%q} =! anaconda.Tweet.User.ScreenName{%q}",
				referenceTweet.UserScreenName,
				tweet.UserScreenName)
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
					"ReferenceTweet.Medias[%d]{%q} "+
						"=! anaconda.Tweet.ExtendedEntities.Media[%d]{%q}",
					j,
					referenceMedia,
					j,
					media)
			}
		}
	}
}
