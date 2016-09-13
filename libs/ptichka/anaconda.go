package ptichka

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

type anacondaTweets []anaconda.Tweet

func (tweets anacondaTweets) toTweets() (tweetsByDate, error) {
	t := make(tweetsByDate, len(tweets))
	for i := range tweets {
		date, err := time.Parse(time.RubyDate, tweets[i].CreatedAt)
		if err != nil {
			return t, err
		}

		m := tweets[i].ExtendedEntities.Media
		medias := make([]*media, 0, len(m))
		for j := range m {
			url, err := url.Parse(m[j].Media_url_https)
			if err == nil {
				medias = append(medias, &media{
					IDStr:         m[j].Id_str,
					MediaURLHttps: url.String()})
			}
		}

		var text string
		if tweets[i].RetweetedStatus == nil {
			text = tweets[i].Text
		} else {
			text = fmt.Sprintf("RT @%s: %s",
				tweets[i].RetweetedStatus.User.ScreenName,
				tweets[i].RetweetedStatus.Text)
		}
		t[i] = &tweet{
			IDStr:          tweets[i].IdStr,
			UserScreenName: tweets[i].User.ScreenName,
			Date:           date,
			Text:           text,
			Medias:         medias}
	}
	return t, nil
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
	tweets, err := api.GetHomeTimeline(nil)
	return tweets, err
}
