package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := loadConfig(".twitmarc.toml.example")
	if err != nil {
		t.Errorf("Error on loadConfig(\".twitmarc.toml.example\"): %v", err)
	}

	if config.CacheFile != ".twitma.json" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").CacheFile == %v, want %v",
			config.CacheFile,
			".twitma.json")
	}
	if config.Label != "[twitter] " {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Label == %v, want %v",
			config.Label,
			"[twitter] ")
	}
	if config.Verbose != true {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Verbose == %v, want %v",
			config.Verbose,
			true)
	}
	if config.Mail.From.Address != "noreply@example.com" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.From.Address == %v, want %v",
			config.Mail.From.Address,
			"noreply@example.com")
	}
	if config.Mail.To.Address != "your.mail@example.org" {
		t.Errorf(
			"loadConfig(\".twitmarc.toml.example\").Mail.To.Address == %v, want %v",
			config.Mail.To.Address,
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
