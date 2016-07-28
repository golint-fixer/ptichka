package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

// AnacondaTweets is a collection of the anaconda.Tweet.
type AnacondaTweets []anaconda.Tweet

func fetchTweets(config *config) (AnacondaTweets, error) {
	anaconda.SetConsumerKey(config.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(config.Twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(
		config.Twitter.AccessToken,
		config.Twitter.AccessTokenSecret)
	anacondaTweets, err := api.GetHomeTimeline(nil)
	return anacondaTweets, err
}

func (anacondaTweets AnacondaTweets) toTweets() TweetsByDate {
	tweets := make(TweetsByDate, len(anacondaTweets))
	for i := range anacondaTweets {
		date, err := time.Parse(time.RubyDate, anacondaTweets[i].CreatedAt)
		ifError(err, "Error on time.Parse CreatedAt: %s")

		var medias []string
		for j := range anacondaTweets[i].ExtendedEntities.Media {
			url, err := url.Parse(
				anacondaTweets[i].ExtendedEntities.Media[j].Media_url_https)
			if err == nil {
				medias = append(medias, url.String())
			}
		}
		var text string
		if anacondaTweets[i].RetweetedStatus == nil {
			text = anacondaTweets[i].Text
		} else {
			text = fmt.Sprintf("RT @%s: %s",
				anacondaTweets[i].RetweetedStatus.User.ScreenName,
				anacondaTweets[i].RetweetedStatus.Text)
		}
		tweets[i] = Tweet{
			IDStr:          anacondaTweets[i].IdStr,
			UserScreenName: anacondaTweets[i].User.ScreenName,
			Date:           date,
			Text:           text,
			Medias:         medias}
	}
	return tweets
}
