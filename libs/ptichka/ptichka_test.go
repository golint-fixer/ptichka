package ptichka

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"testing"
	"text/template"
	"time"
)

func TestTweetsByDateSrot(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	got := tweetsByDate{
		&tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		&tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)},
		&tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)}}

	wont := tweetsByDate{
		&tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)},
		&tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		&tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)}}

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

func dummyPtichka(fetcher Fetcher, sender Sender) []error {
	var err error
	var errs []error

	cacheFile1, err := ioutil.TempFile(os.TempDir(), "ptichka_cache1")
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	defer func() { _ = os.Remove(cacheFile1.Name()) }()

	cacheFile2, err := ioutil.TempFile(os.TempDir(), "ptichka_cache2")
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	defer func() { _ = os.Remove(cacheFile2.Name()) }()

	emptyJSON := []byte("[]")
	for _, f := range []*os.File{cacheFile1, cacheFile2} {
		if _, err := f.Write(emptyJSON); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	logFile1, err := ioutil.TempFile(os.TempDir(), "ptichka_log1")
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	defer func() { _ = os.Remove(logFile1.Name()) }()

	logFile2, err := ioutil.TempFile(os.TempDir(), "ptichka_log2")
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	defer func() { _ = os.Remove(logFile2.Name()) }()

	tmpl, err := template.New("config").Parse(`
[[accounts]]

  cache_file = "{{.CacheFile1}}"
  label = "[twitter1] "
  log_file = "{{.LogFile1}}"

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

  cache_file = "{{.CacheFile2}}"
  label = "[twitter2] "
  log_file = "{{.LogFile2}}"

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
`)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	var configBuffer bytes.Buffer

	err = tmpl.Execute(&configBuffer, struct {
		CacheFile1 string
		LogFile1   string
		CacheFile2 string
		LogFile2   string
	}{
		CacheFile1: cacheFile1.Name(),
		LogFile1:   logFile1.Name(),
		CacheFile2: cacheFile2.Name(),
		LogFile2:   logFile2.Name()})

	if err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(
		errs,
		Ptichka("", configBuffer.String(), fetcher, sender)...)

	return errs
}

// TestPtichka is an integration test.
func TestPtichka(t *testing.T) {
	errs := dummyPtichka(&dummyFetcher{}, &dummySMTPSender{})
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}
