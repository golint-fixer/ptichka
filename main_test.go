package main

import (
	"testing"
)

func TestTweetBody(t *testing.T) {
	got, err := tweetBody(tweet{"1234", "johndoe", "foo bar"})
	if err != nil {
		t.Errorf("Error on tweetBody(tweet{...}): %v", err)
	}

	wont := `@johndoe

foo bar

https://twitter.com/johndoe/status/1234`

	if got != wont {
		t.Errorf(`tweetBody("bar")
get:
%v

wont:
%v`, got, wont)
	}
}
