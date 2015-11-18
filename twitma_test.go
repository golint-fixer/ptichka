package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := loadConfig(".twitmarc.toml.example")
	if err != nil {
		t.Errorf("Error on loadConfig(\".twitmarc.toml.example\"): %v", err)
	}

	if config.CacheFile != ".twitma.json" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Count == %v, want %v",
			config.CacheFile,
			".twitma.json")
	}
	if config.Count != 200 {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Count == %v, want %v",
			config.Count,
			200)
	}
	if config.Label != "twitter" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Count == %v, want %v",
			config.Count,
			201)
	}
	if config.Verbose != true {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\"). == %v, want %v",
			config.Verbose,
			true)
	}
	if config.Mail.From != "noreply@example.com" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.From == %v, want %v",
			config.Mail.From,
			"noreply@example.com")
	}
	if config.Mail.To != "your.mail@example.org" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.To == %v, want %v",
			config.Mail.To,
			"your.mail@example.org")
	}
	if config.Mail.SMTP.Address != "mail.example.com" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.SMTP.Address == %v, want %v",
			config.Mail.SMTP.Address,
			"mail.example.com")
	}
	if config.Mail.SMTP.Authentication != "plain" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.SMTP.Authentication == %v, want %v",
			config.Mail.SMTP.Authentication,
			"plain")
	}
	if config.Mail.SMTP.Password != "your_password" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.SMTP.Password == %v, want %v",
			config.Mail.SMTP.Password,
			"your_password")
	}
	if config.Mail.SMTP.Port != 25 {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.SMTP.Port == %v, want %v",
			config.Mail.SMTP.Port,
			25)
	}
	if config.Mail.SMTP.SSL != false {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.SMTP.SSL == %v, want %v",
			config.Mail.SMTP.SSL,
			false)
	}
	if config.Mail.SMTP.TLS != false {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.SMTP.TLS == %v, want %v",
			config.Mail.SMTP.TLS,
			false)
	}
	if config.Mail.SMTP.UserName != "your.mail@example.org" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.SMTP.UserName == %v, want %v",
			config.Mail.SMTP.UserName,
			"your.mail@example.org")
	}
	if config.Twitter.ConsumerKey != "your-consumer-key" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Twitter.ConsumerKey == %v, want %v",
			config.Twitter.ConsumerKey,
			"your-consumer-key")
	}
	if config.Twitter.ConsumerSecret != "your-consumer-secret" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Twitter.ConsumerSecret == %v, want %v",
			config.Twitter.ConsumerSecret,
			"your-consumer-secret")
	}
	if config.Twitter.AccessToken != "your-access-token" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Twitter.AccessToken == %v, want %v",
			config.Twitter.AccessToken,
			"your-access-token")
	}
	if config.Twitter.AccessTokenSecret != "your-access-token-secret" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Twitter.AccessTokenSecret == %v, want %v",
			config.Twitter.AccessTokenSecret,
			"your-access-token-secret")
	}
}

func TestLoadIds(t *testing.T) {
	var jsonBlob = []byte(`[
"first-id",
"second-id"
]`)

	tempFile, err := ioutil.TempFile(os.TempDir(), "twitma_test_ids_")
	if err != nil {
		t.Errorf("Error on create temporary file: %v", err)
	}

	ioutil.WriteFile(tempFile.Name(), jsonBlob, 0644)

	ids, err := loadIds(tempFile.Name())
	if err != nil {
		t.Errorf("Error on loadIds(tempFile): %v", err)
	}

	defer os.Remove(tempFile.Name())

	var got, wont string

	got = ids[0]
	wont = "first-id"
	if got != wont {
		t.Errorf("loadIds(jsonBlob)[n] == %v, want %v", got, wont)
	}

	got = ids[1]
	wont = "second-id"
	if got != wont {
		t.Errorf("loadIds(jsonBlob)[n] == %v, want %v", got, wont)
	}
}
