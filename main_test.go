package main

import (
	"testing"
)

func TestTweetBody(t *testing.T) {
	got, err := tweetBody(tweet{ID: "1234", User: "johndoe", Text: "foo bar"})
	if err != nil {
		t.Errorf("Error on tweetBody(tweet{...}): %v", err)
	}

	wont := `@johndoe

foo bar

https://twitter.com/johndoe/status/1234`

	if got != wont {
		t.Errorf(`tweetBody(tweet{ID: "1234", User: "johndoe", Text: "foo bar"})
get:
%v

wont:
%v`, got, wont)
	}
}
