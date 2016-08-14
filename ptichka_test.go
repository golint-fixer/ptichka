package ptichka

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestTweetsByDateSrot(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	got := tweetsByDate{
		tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)},
		tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)}}

	wont := tweetsByDate{
		tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)},
		tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)}}

	sort.Sort(got)

	if !reflect.DeepEqual(got, wont) {
		t.Errorf("%v != %v", got, wont)
	}
}

func TestTweetText(t *testing.T) {
	got, err := tweet{
		IDStr:          "1234",
		UserScreenName: "johndoe",
		Text:           "Hello &amp; world!"}.text()
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

type dummyFetcher struct{}

func (f *dummyFetcher) fetch(config *configuration) (anacondaTweets, error) {
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

	return anacondaTweets, err
}

type dummySMTPSender struct{}

func (e *dummySMTPSender) send(config *configuration, msg []byte) error {
	return nil
}

// TestPtichka is an integration test.
func TestPtichka(t *testing.T) {
	cacheFile1, err := ioutil.TempFile(os.TempDir(), "ptichka_cache1")
	defer func() { _ = os.Remove(cacheFile1.Name()) }()
	if err != nil {
		t.Error(err)
	}
	cacheFile2, err := ioutil.TempFile(os.TempDir(), "ptichka_cache2")
	defer func() { _ = os.Remove(cacheFile2.Name()) }()
	if err != nil {
		t.Error(err)
	}
	emptyJSON := []byte("[]")
	for _, f := range []*os.File{cacheFile1, cacheFile2} {
		if _, err := f.Write(emptyJSON); err != nil {
			t.Error(err)
		}
	}

	logFile1, err := ioutil.TempFile(os.TempDir(), "ptichka_log1")
	defer func() { _ = os.Remove(logFile1.Name()) }()
	if err != nil {
		t.Error(err)
	}
	logFile2, err := ioutil.TempFile(os.TempDir(), "ptichka_log2")
	defer func() { _ = os.Remove(logFile2.Name()) }()
	if err != nil {
		t.Error(err)
	}

	var rawConfig = fmt.Sprintf(`
[[accounts]]

  cache_file = "%s"
  label = "[twitter1] "
  log_file = "%s"

  [accounts.twitter]
    consumer_key = "foo1"
    consumer_secret = "bar1"
    access_token = "baz1"
    access_token_secret = "xyz1"

  [accounts.mail]

    [accounts.mail.to]
      address = "to1@example.org"
      name = "Test recipient 1"

    [accounts.mail.from]
      address = "from1@example.org"
      name = "Test sender 1"

    [accounts.mail.smtp]
      address = "smpt.example.org"
      password = "password"
      port = 25
      user_name = "to1@example.org"

[[accounts]]

  cache_file = "%s"
  label = "[twitter2] "
  log_file = "%s"

  [accounts.twitter]
    consumer_key = "foo2"
    consumer_secret = "bar2"
    access_token = "baz2"
    access_token_secret = "xyz2"

  [accounts.mail]

    [accounts.mail.to]
      address = "to2@example.org"
      name = "Test recipient 2"

    [accounts.mail.from]
      address = "from2@example.org"
      name = "Test sender 2"

    [accounts.mail.smtp]
      address = "smpt.example.org"
      password = "password"
      port = 25
      user_name = "to2@example.org"

`, cacheFile1.Name(), logFile1.Name(), cacheFile2.Name(), logFile2.Name())

	errs := Ptichka("", rawConfig, &dummyFetcher{}, &dummySMTPSender{})
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}
