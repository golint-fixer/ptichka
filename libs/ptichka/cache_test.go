package ptichka

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadCache(t *testing.T) {
	var jsonBlob = []byte(`["first-id", "second-id"]`)

	tempFile, err := ioutil.TempFile(os.TempDir(), "ptichka_test_ids_")
	if err != nil {
		t.Errorf("Error on create temporary file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()

	tempFileName := tempFile.Name()
	if err := ioutil.WriteFile(tempFileName, jsonBlob, 0644); err != nil {
		t.Fatalf("WriteFile %s: %v", tempFileName, err)
	}

	ids, err := loadCache(tempFile.Name())
	if err != nil {
		t.Errorf("Error on loadCache(tempFile): %v", err)
	}

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
