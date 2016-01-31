package main

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"text/template"
)

// filterTweets reject tweets thats intersect
// with "cached" tweets's oldIds (thats processed in previous time).
// It returns the slice of the newIds thats was not processed in prevented time.
func filterTweets(
	tweets TweetsByDate,
	oldIds []string,
	config *config) []string {
	newIds := oldIds

	for _, currentTweet := range tweets {
		if !contains(oldIds, currentTweet.ID) {
			newIds = append(newIds, currentTweet.ID)

			// for example "[twitter] @JohnDoe 1970-01-01 00:00 +0000"
			subject := fmt.Sprintf(
				"%s@%s %s",
				config.Label,
				currentTweet.User,
				currentTweet.Date.Format("2006-01-02 15:04 -0700"))

			if config.Verbose {
				print("Sending: " + subject)
			}

			body, err := tweetBody(Tweet{
				ID:   currentTweet.ID,
				User: currentTweet.User,
				Text: currentTweet.Text})
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
				defer response.Body.Close()

				tempDir, err := ioutil.TempDir(
					os.TempDir(),
					fmt.Sprintf("tgtm_%s", currentTweet.ID))
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(tempDir)

				_, fileName := filepath.Split(media)

				tempFilePath := fmt.Sprintf("%s/%s", tempDir, fileName)

				tempFile, err := os.Create(tempFilePath)
				defer tempFile.Close()

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
	return newIds
}

func tweetBody(t Tweet) (string, error) {
	tmpl, err := template.New("tweet").Parse(
		`@{{.User}}

{{.Text}}

https://twitter.com/{{.User}}/status/{{.ID}}`)
	if err != nil {
		panic(err)
	}

	var x bytes.Buffer

	err = tmpl.Execute(&x, Tweet{
		ID:   t.ID,
		User: t.User,
		Text: html.UnescapeString(t.Text)})

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
