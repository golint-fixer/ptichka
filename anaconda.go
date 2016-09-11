package ptichka

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type anacondaTweets []anaconda.Tweet

func (anacondaTweets anacondaTweets) toTweets() (tweetsByDate, error) {
	tweets := make(tweetsByDate, len(anacondaTweets))
	for i := range anacondaTweets {
		date, err := time.Parse(time.RubyDate, anacondaTweets[i].CreatedAt)
		if err != nil {
			return tweets, err
		}

		m := anacondaTweets[i].ExtendedEntities.Media
		medias := make([]media, 0, len(m))
		for j := range m {
			url, err := url.Parse(m[j].Media_url_https)
			if err == nil {
				medias = append(medias, media{
					IDStr:         m[j].Id_str,
					MediaURLHttps: url.String()})
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
		tweets[i] = tweet{
			IDStr:          anacondaTweets[i].IdStr,
			UserScreenName: anacondaTweets[i].User.ScreenName,
			Date:           date,
			Text:           text,
			Medias:         medias}
	}
	return tweets, nil
}

// Fetcher is an interface with fetch method.
type Fetcher interface {
	fetch(config *configuration) (anacondaTweets, error)
}

// AnacondaFetcher has a method to download tweets
// realized through the Anaconda package.
type AnacondaFetcher struct{}

// fetch method download Twitter home timeline.
func (f *AnacondaFetcher) fetch(config *configuration) (anacondaTweets, error) {
	anaconda.SetConsumerKey(config.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(config.Twitter.ConsumerSecret)
	api := anaconda.NewTwitterApi(
		config.Twitter.AccessToken,
		config.Twitter.AccessTokenSecret)
	anacondaTweets, err := api.GetHomeTimeline(nil)
	return anacondaTweets, err
}
