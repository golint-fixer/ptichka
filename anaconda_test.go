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
        {
          "id_str": "foo",
          "Media_url_https": "https://pbs.twimg.com/media/foo.jpg"
        }, {
          "id_str": "should_be_rejected",
          "Media_url_https": ":this://is.wrong/url/should/be/rejected"
        }, {
          "id_str": "bar",
          "Media_url_https": "https://pbs.twimg.com/media/bar.png"
        }
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
        {
          "id_str": "xyz",
          "Media_url_https": "https://pbs.twimg.com/media/xyz.gif"
        }
      ]
    }
  }
]`)

	var anacondaTweets anacondaTweets
	err := json.Unmarshal(jsonBlob, &anacondaTweets)
	if err != nil {
		log.Fatal(err)
	}
	ptichkaTweets, err := anacondaTweets.toTweets()
	if err != nil {
		log.Fatal(err)
	}

	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	referenceTweets := tweetsByDate{
		tweet{
			IDStr:          "111111111111111111",
			Date:           time.Date(1970, 1, 1, 1, 0, 0, 0, utc),
			UserScreenName: "johndoe",
			Text:           "RT @ivanivanov: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim"},
		tweet{
			IDStr:          "222222222222222222",
			Date:           time.Date(1970, 1, 2, 2, 0, 0, 0, utc),
			UserScreenName: "johndoe",
			Text:           "foo bar",
			Medias: []media{
				{IDStr: "foo", MediaURLHttps: "https://pbs.twimg.com/media/foo.jpg"},
				{IDStr: "bar", MediaURLHttps: "https://pbs.twimg.com/media/bar.png"}}},

		tweet{
			IDStr:          "333333333333333333",
			Date:           time.Date(1970, 1, 3, 3, 0, 0, 0, utc),
			UserScreenName: "gunterschmidt",
			Text:           "baz xyz",
			Medias: []media{
				{IDStr: "xyz", MediaURLHttps: "https://pbs.twimg.com/media/xyz.gif"}}}}

	for i := range referenceTweets {
		referenceTweet := referenceTweets[i]
		ptichkaTweet := ptichkaTweets[i]

		if referenceTweet.IDStr != ptichkaTweet.IDStr {
			t.Errorf(
				"ReferenceTweet.IDStr{%q} =! anaconda.Tweet.IdStr{%q}",
				referenceTweet.IDStr,
				ptichkaTweet.IDStr)
		}

		if !referenceTweet.Date.Equal(ptichkaTweet.Date) {
			t.Errorf(
				"ReferenceTweet.Date{%q} =! anaconda.Tweet.createdAt{%q}",
				referenceTweet.Date,
				ptichkaTweet.Date)
		}

		if referenceTweet.UserScreenName != ptichkaTweet.UserScreenName {
			t.Errorf(
				"ReferenceTweet.UserScreenName{%q} =! anaconda.Tweet.User.ScreenName{%q}",
				referenceTweet.UserScreenName,
				ptichkaTweet.UserScreenName)
		}

		if referenceTweet.Text != ptichkaTweet.Text {
			t.Errorf(
				"ReferenceTweet.Text{%q} =! anaconda.Tweet.Text{%q}",
				referenceTweet.Text,
				ptichkaTweet.Text)
		}

		for j := range referenceTweet.Medias {
			referenceMedia := referenceTweet.Medias[j]
			ptichkaMedia := ptichkaTweet.Medias[j]

			if referenceMedia != ptichkaMedia {
				t.Errorf(
					"ReferenceTweet.Medias[%[1]d]%q "+
						"=! anaconda.Tweet.ExtendedEntities.Media[%[1]d]%q",
					j,
					referenceMedia,
					ptichkaMedia)
			}
		}
	}
}
