// Package ptichka fetch timeline tweets and sends by SMTP.
package ptichka

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"gopkg.in/jordan-wright/email.v2"
)

// Version is an package version.
const Version = "0.6.9"

// tweet is a simplified anaconda.Tweet.
type tweet struct {
	IDStr          string
	UserScreenName string
	Date           time.Time
	Text           string
	Medias         []media
}

// title returns string with label, user name and time
// for example "[twitter] @JohnDoe 1970-01-01 00:00 +0000"
func (t tweet) title(config *config) string {
	return fmt.Sprintf(
		"%s@%s %s",
		config.Label,
		t.UserScreenName,
		t.Date.Format("2006-01-02 15:04 -0700"))
}

func (t tweet) text() (string, error) {
	tmpl, err := template.New("tweet").Parse(
		`@{{.UserScreenName}}

{{.Text}}

https://twitter.com/{{.UserScreenName}}/status/{{.IDStr}}`)
	if err != nil {
		return "", err
	}

	var x bytes.Buffer

	err = tmpl.Execute(&x, tweet{
		IDStr:          t.IDStr,
		UserScreenName: t.UserScreenName,
		Text:           html.UnescapeString(t.Text)})

	return x.String(), err
}

// media is an simplified anaconda EntityMedia.
type media struct {
	IDStr          string
	MediaURLHttps  string
	DownloadedPath string
}

// tweetsByDate is a slice of tweet
// with ability to be sorted by date from older to newer.
type tweetsByDate []tweet

// <https://github.com/wskinner/anaconda/commit/d0c12d8fba671d7d5ce27d3abd1809aedcc59195>,
// <http://nerdyworm.com/blog/2013/05/15/sorting-a-slice-of-structs-in-go/>.

// Len is the number of elements in the collection.
func (a tweetsByDate) Len() int { return len(a) }

// Less reports whether the element with index i should sort before
// the element with index j.
func (a tweetsByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Swap swaps the elements with indexes i and j.
func (a tweetsByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// Fly fetch home timeline and sends by SMTP.
func Fly(config *config, errCh chan<- error) {
	oldIds, err := loadCache(config.CacheFile)
	if err != nil {
		errCh <- err
		return
	}

	messages, err := fetch(config, oldIds)
	if err != nil {
		errCh <- err
		return
	}

	var newIds []string

	// Sends mails without goroutines
	// to preserve sort order of the mails (tweets).
	for IDStr, message := range messages {
		err = smtp.SendMail(
			fmt.Sprintf("%s:%d", config.Mail.SMTP.Address, config.Mail.SMTP.Port),
			smtp.PlainAuth(
				"",
				config.Mail.SMTP.UserName,
				config.Mail.SMTP.Password,
				config.Mail.SMTP.Address),
			config.mailFrom(),
			[]string{config.mailTo()},
			message)
		if err != nil {
			errCh <- err
			break
		}
		newIds = append(newIds, IDStr)
	}

	err = saveCache(config.CacheFile, append(oldIds, newIds...))
	if err != nil {
		errCh <- err
	}
}

// fetch fetch home timeline and returns
// array of the raw email messages
// (with headers, text and attachments encoded in base64).
func fetch(config *config, oldIds []string) (map[string][]byte, error) {
	messages := make(map[string][]byte)

	anacondaTweets, err := fetchTweets(config)
	if err != nil {
		return messages, err
	}

	tweets, err := anacondaTweets.toTweets()
	if err != nil {
		return messages, err
	}

	sort.Sort(tweets)

	var newTweets tweetsByDate
	newMedias := make(map[string]media)

	tempDirPath, err := ioutil.TempDir(os.TempDir(), "ptichka_")
	defer func() { _ = os.RemoveAll(tempDirPath) }()
	if err != nil {
		return messages, err
	}

	mediaCh := make(chan media)
	mediaErrCh := make(chan error)
	for _, currentTweet := range tweets {
		if contains(oldIds, currentTweet.IDStr) {
			continue
		}

		newTweets = append(newTweets, currentTweet)

		for _, m := range currentTweet.Medias {
			newMedias[m.IDStr] = m
		}
	}

	for _, m := range newMedias {
		go getMedia(m, tempDirPath, mediaCh, mediaErrCh)
	}

	for range newMedias {
		newMedia := <-mediaCh
		newMedias[newMedia.IDStr] = newMedia
	}

	var errs []error
	for range newMedias {
		err = <-mediaErrCh
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return messages, err
	}

	for _, newTweet := range newTweets {
		message := email.NewEmail()
		message.From = config.mailFrom()
		message.To = []string{config.mailTo()}
		message.Subject = newTweet.title(config)

		text, err := newTweet.text()
		if err != nil {
			return messages, err
		}
		message.Text = []byte(text)

		for _, m := range newTweet.Medias {
			// Follow the assumption that the medias id_str is unique
			// for the whole Twitter.
			_, err = message.AttachFile(newMedias[m.IDStr].DownloadedPath)
			if err != nil {
				return messages, err
			}
		}

		rawMessage, err := message.Bytes()
		if err != nil {
			return messages, err
		}
		messages[newTweet.IDStr] = rawMessage
	}

	return messages, nil
}

func getMedia(
	newMedia media,
	tempDirPath string,
	ch chan<- media,
	errCh chan<- error) {

	response, err := http.Get(newMedia.MediaURLHttps)
	if err != nil {
		ch <- newMedia
		errCh <- err
		return
	}
	defer func() { _ = response.Body.Close() }()

	_, fileName := filepath.Split(newMedia.MediaURLHttps)

	tempFilePath := fmt.Sprintf("%s/%s", tempDirPath, fileName)

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		ch <- newMedia
		errCh <- err
		return
	}
	defer func() { _ = tempFile.Close() }()

	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		ch <- newMedia
		errCh <- err
		return
	}

	newMedia.DownloadedPath = tempFilePath

	ch <- newMedia
	errCh <- nil
}

func contains(ids []string, id string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}
