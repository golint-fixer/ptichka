package main

import (
	"testing"
)

func TestTweetBody(t *testing.T) {
	got, err := tweetBody(Tweet{
		IDStr:          "1234",
		UserScreenName: "johndoe",
		Text:           "Hello &amp; world!"})

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
