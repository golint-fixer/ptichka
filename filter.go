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
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/gomail.v2"
)

// filterTweets reject tweets thats intersect
// with "cached" tweets's oldIds (thats processed in previous time).
// It returns the slice of the newIds which was not processed in previous time.
func filterTweets(
	tweets TweetsByDate,
	oldIds []string,
	config *config) []string {
	newIds := oldIds

	for _, currentTweet := range tweets {
		if !contains(oldIds, currentTweet.IDStr) {
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

			message := gomail.NewMessage()
			message.SetHeader("From", from.String())
			message.SetHeader("To", to.String())
			message.SetHeader("Subject", subject)
			message.SetBody("text/plain", body)

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

				message.Attach(tempFilePath)
			}

			dialer := gomail.NewPlainDialer(
				config.Mail.SMTP.Address,
				config.Mail.SMTP.Port,
				config.Mail.SMTP.UserName,
				config.Mail.SMTP.Password)

			if err := dialer.DialAndSend(message); err != nil {
				panic(err)
			}

			if config.Verbose {
				print("\n")
			}
		}
	}
	return newIds //FIXME: it's actually all ids
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

	return x.String(), err
}

func contains(ids []string, id string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}
