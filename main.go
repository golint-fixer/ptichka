package main

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"gopkg.in/jordan-wright/email.v2"
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
	config, err := loadConfig(os.Args[1])
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

	newIds := oldIds

	for _, currentTweet := range tweets {
		if contains(oldIds, currentTweet.IDStr) {
			continue
		}

		newIds = append(newIds, currentTweet.IDStr)

		// for example "[twitter] @JohnDoe 1970-01-01 00:00 +0000"
		subject := fmt.Sprintf(
			"%s@%s %s",
			config.Label,
			currentTweet.UserScreenName,
			currentTweet.Date.Format("2006-01-02 15:04 -0700"))

		if config.Verbose {
			print("Sending: " + subject)
		}

		body, err := tweetBody(Tweet{
			IDStr:          currentTweet.IDStr,
			UserScreenName: currentTweet.UserScreenName,
			Text:           currentTweet.Text})
		if err != nil {
			panic(err)
		}

		from := mail.Address{
			Name:    config.Mail.From.Name,
			Address: config.Mail.From.Address}
		to := mail.Address{
			Name:    config.Mail.To.Name,
			Address: config.Mail.To.Address}

		message := email.NewEmail()
		message.From = from.String()
		message.To = []string{to.String()}
		message.Subject = subject
		message.Text = []byte(body)

		for _, media := range currentTweet.Medias {
			response, err := http.Get(media)
			if err != nil {
				log.Fatal(err)
			}
			defer func() { _ = response.Body.Close() }()

			tempDir, err := ioutil.TempDir(
				os.TempDir(),
				fmt.Sprintf("ptichka_%s", currentTweet.IDStr))
			if err != nil {
				log.Fatal(err)
			}
			defer func() { _ = os.Remove(tempDir) }()

			_, fileName := filepath.Split(media)

			tempFilePath := fmt.Sprintf("%s/%s", tempDir, fileName)

			tempFile, err := os.Create(tempFilePath)
			defer func() { _ = tempFile.Close() }()

			_, err = io.Copy(tempFile, response.Body)
			if err != nil {
				log.Fatal(err)
			}

			_, err = message.AttachFile(tempFilePath)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = message.Send(
			fmt.Sprintf("%s:%d", config.Mail.SMTP.Address, config.Mail.SMTP.Port),
			smtp.PlainAuth(
				"",
				config.Mail.SMTP.UserName,
				config.Mail.SMTP.Password,
				config.Mail.SMTP.Address))
		if err != nil {
			panic(err)
		}

		if config.Verbose {
			print("\n")
		}
	}

	if err := saveCache(config.CacheFile, newIds); err != nil {
		log.Fatal(err)
	}
}

func tweetBody(t Tweet) (string, error) {
	tmpl, err := template.New("tweet").Parse(
		`@{{.UserScreenName}}

{{.Text}}

https://twitter.com/{{.UserScreenName}}/status/{{.IDStr}}`)
	if err != nil {
		panic(err)
	}

	var x bytes.Buffer

	err = tmpl.Execute(&x, Tweet{
		IDStr:          t.IDStr,
		UserScreenName: t.UserScreenName,
		Text:           html.UnescapeString(t.Text)})

	return x.String(), err // TODO: remove extra conversion byte->string->byte
}

func contains(ids []string, id string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}
