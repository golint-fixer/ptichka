package main

import (
	"testing"
)

func TestTweetBody(t *testing.T) {
	got, err := tweetBody(tweet{
		ID:   "1234",
		User: "johndoe",
		Text: "Hello &amp; world!"})

	if err != nil {
		t.Errorf("Error on tweetBody(tweet{...}): %v", err)
	}

	wont := `@johndoe

Hello & world!

https://twitter.com/johndoe/status/1234`

	if got != wont {
		t.Errorf(`tweetBody(tweet{ID: "1234", User: "johndoe", Text: "Hello &amp; world!"})
get:
%v

wont:
%v`, got, wont)
	}
}
