package main

import (
	"log"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}

	got := TweetsByDate{
		Tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)}}

	wont := TweetsByDate{
		Tweet{Date: time.Date(1970, 3, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(2000, 2, 1, 1, 0, 0, 0, utc)},
		Tweet{Date: time.Date(2123, 1, 1, 1, 0, 0, 0, utc)}}

	sort.Sort(got)

	if !reflect.DeepEqual(got, wont) {
		t.Errorf("%v != %v", got, wont)
	}
}
