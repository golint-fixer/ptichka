package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadCache(t *testing.T) {
	var jsonBlob = []byte(`["first-id", "second-id"]`)

	tempFile, err := ioutil.TempFile(os.TempDir(), "tgtm_test_ids_")
	if err != nil {
		t.Errorf("Error on create temporary file: %v", err)
	}

	ioutil.WriteFile(tempFile.Name(), jsonBlob, 0644)

	ids, err := loadCache(tempFile.Name())
	if err != nil {
		t.Errorf("Error on loadCache(tempFile): %v", err)
	}

	defer os.Remove(tempFile.Name())

	var got, wont string

	got = ids[0]
	wont = "first-id"
	if got != wont {
		t.Errorf("loadCache(jsonBlob)[n] == %v, want %v", got, wont)
	}

	got = ids[1]
	wont = "second-id"
	if got != wont {
		t.Errorf("loadCache(jsonBlob)[n] == %v, want %v", got, wont)
	}
}
