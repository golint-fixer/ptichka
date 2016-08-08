// Package ptichka fetch timeline tweets and sends by SMTP.
package ptichka

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
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
const Version = "0.6.18"

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
func (t tweet) title(config *configuration) string {
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
func Fly(
	config *configuration,
	errCh chan<- error,
	infLogger, errLogger *log.Logger) {

	oldIds, err := loadCache(config.CacheFile)
	if err != nil {
		errCh <- err
		return
	}

	infLogger.Printf("%sCache loaded: %s", config.Label, config.CacheFile)

	messages, err := fetch(config, oldIds, infLogger, errLogger)
	if err != nil {
		errCh <- err
		return
	}

	infLogger.Printf("%sFetched messages: %d", config.Label, len(messages))

	var newIds []string

	// Sends mails without goroutines
	// to preserve sort order of the mails (tweets).
	var smtpErr error

	for IDStr, message := range messages {
		smtpErr = smtp.SendMail(
			fmt.Sprintf("%s:%d", config.Mail.SMTP.Address, config.Mail.SMTP.Port),
			smtp.PlainAuth(
				"",
				config.Mail.SMTP.UserName,
				config.Mail.SMTP.Password,
				config.Mail.SMTP.Address),
			config.mailFrom(),
			[]string{config.mailTo()},
			message)

		if smtpErr != nil {
			break
		}

		infLogger.Printf("%sMessages sent: %s", config.Label, IDStr)
		newIds = append(newIds, IDStr)
	}

	var cacheErr error
	if len(newIds) > 0 {
		cacheErr = saveCache(config.CacheFile, append(oldIds, newIds...))
	}

	if smtpErr != nil {
		errLogger.Printf("%sSent messages: %d, messages does not sent: %d",
			config.Label, len(newIds), len(messages)-len(newIds))

		errCh <- smtpErr
		return
	}
	infLogger.Printf("%sSent messages: %d", config.Label, len(newIds))

	if cacheErr != nil {
		errLogger.Printf("%sCache does not saved: %s",
			config.Label, config.CacheFile)

		errCh <- cacheErr
		return
	}
	if len(newIds) > 0 {
		infLogger.Printf("%sCache saved: %s", config.Label, config.CacheFile)
	}

	errCh <- nil
}

// fetch fetch home timeline and returns
// array of the raw email messages
// (with headers, text and attachments encoded in base64).
func fetch(
	config *configuration,
	oldIds []string,
	infLogger, errLogger *log.Logger) (map[string][]byte, error) {

	messages := make(map[string][]byte)

	anacondaTweets, err := fetchTweets(config)
	if err != nil {
		errLogger.Printf("%sTweets does not fetched", config.Label)
		return messages, err
	}
	infLogger.Printf("%sTweets fetched: %d", config.Label, len(anacondaTweets))

	tweets, err := anacondaTweets.toTweets()
	if err != nil {
		return messages, err
	}

	sort.Sort(tweets)

	var newTweets tweetsByDate
	newMedias := make(map[string]media)

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

	if len(newTweets) > 0 {
		infLogger.Printf("%sNew tweets fetched: %d", config.Label, len(newTweets))
	}

	var tempDirPath string
	if len(newMedias) > 0 {
		tempDirPath, err = ioutil.TempDir(os.TempDir(), "ptichka_")
		defer func() {
			err := os.RemoveAll(tempDirPath)
			if err != nil {
				errLogger.Printf("%sTemporary directory does not removed: %s %v",
					config.Label, tempDirPath, err)
			} else {
				infLogger.Printf("%sTemporary directory removed: %s",
					config.Label, tempDirPath)
			}
		}()
		if err != nil {
			errLogger.Printf("%sTemporary directory does not created: %s",
				config.Label, tempDirPath)
			return messages, err
		}
		infLogger.Printf("%sTemporary directory created: %s",
			config.Label, tempDirPath)
	}
	for _, m := range newMedias {
		go getMedia(
			config,
			m,
			tempDirPath,
			mediaCh, mediaErrCh,
			infLogger, errLogger)
	}

	for range newMedias {
		newMedia := <-mediaCh
		newMedias[newMedia.IDStr] = newMedia
	}
	var errs []error
	for range newMedias {
		err := <-mediaErrCh
		if err != nil {
			errs = append(errs, err)
		}
	}
	infLogger.Printf("%sAttachments downloaded: %d", config.Label, len(newMedias))
	if len(errs) > 0 {
		errLogger.Printf("%sErrors when downloading attachments: %v",
			config.Label, errs)
		return messages, errs[0]
	}

	for _, newTweet := range newTweets {
		text, err := newTweet.text()
		if err != nil {
			return messages, err
		}

		message := &email.Email{
			From:    config.mailFrom(),
			To:      []string{config.mailTo()},
			Subject: newTweet.title(config),
			Text:    []byte(text)}

		for _, m := range newTweet.Medias {
			// Follow the assumption that the medias id_str is unique
			// for the whole Twitter.
			_, err := message.AttachFile(newMedias[m.IDStr].DownloadedPath)
			if err != nil {
				errLogger.Printf("%sAttaching error: %s %v",
					config.Label, newMedias[m.IDStr].DownloadedPath, err)
				return messages, err
			}
			infLogger.Printf("%sAttached: %s",
				config.Label, newMedias[m.IDStr].DownloadedPath)
		}

		rawMessage, err := message.Bytes()
		if err != nil {
			errLogger.Printf("%sError when generating raw mail: %s %v",
				config.Label, newTweet.IDStr, errs)

			return messages, err
		}

		infLogger.Printf("%sRaw mail generated: %s", config.Label, newTweet.IDStr)
		messages[newTweet.IDStr] = rawMessage
	}

	return messages, nil
}

func getMedia(
	config *configuration,
	newMedia media,
	tempDirPath string,
	ch chan<- media,
	errCh chan<- error,
	infLogger, errLogger *log.Logger) {

	response, err := http.Get(newMedia.MediaURLHttps)
	if err != nil {
		errLogger.Printf("%sError downloading attachment: %s %s %v",
			config.Label, newMedia.IDStr, newMedia.MediaURLHttps, err)

		ch <- newMedia
		errCh <- err
		return
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			errLogger.Printf("%sAttachment response body does not closed: %s %s %v",
				config.Label, newMedia.IDStr, newMedia.MediaURLHttps, err)
		} else {
			infLogger.Printf("%sAttachment response body closed: %s %s",
				config.Label, newMedia.IDStr, newMedia.MediaURLHttps)
		}
	}()
	infLogger.Printf("%sAttachment downloaded: %s %s",
		config.Label, newMedia.IDStr, newMedia.MediaURLHttps)

	_, fileName := filepath.Split(newMedia.MediaURLHttps)

	tempFilePath := fmt.Sprintf("%s/%s", tempDirPath, fileName)

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		errLogger.Printf("%sError creating temporary file: %s %v",
			config.Label, tempDirPath, err)

		ch <- newMedia
		errCh <- err
		return
	}
	defer func() {
		err := tempFile.Close()
		if err != nil {
			errLogger.Printf("%sTemporary file does not closed: %s %v",
				config.Label, tempDirPath, err)
		} else {
			infLogger.Printf("%sTemporary file closed: %s",
				config.Label, tempFilePath)
		}
	}()
	infLogger.Printf("%sTemporary file created: %s", config.Label, tempFilePath)

	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		errLogger.Printf("%sError copying temporary file: %s %v",
			config.Label, tempDirPath, err)

		ch <- newMedia
		errCh <- err
		return
	}
	infLogger.Printf("%sTemporary file copied: %s", config.Label, tempFilePath)

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
