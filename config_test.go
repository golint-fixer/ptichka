package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configs, err := loadConfig(".ptichkarc.toml.example")
	if err != nil {
		t.Errorf("Error on loadConfig(\".ptichkarc.toml.example\"): %v", err)
	}

	config := configs.Accounts[0]

	if config.CacheFile != ".ptichka1.json" {
		t.Errorf("CacheFile == %v, want %v", config.CacheFile, ".ptichka1.json")
	}
	if config.Label != "[twitter account 1] " {
		t.Errorf("Label == %v, want %v", config.Label, "[twitter account 1] ")
	}
	if config.Mail.From.Address != "noreply@example.com" {
		t.Errorf(
			"Mail.From.Address == %v, want %v",
			config.Mail.From.Address,
			"noreply@example.com")
	}
	if config.Mail.To.Address != "your.mail+twitter1@example.org" {
		t.Errorf(
			"Mail.To.Address == %v, want %v",
			config.Mail.To.Address,
			"your.mail+twitter1@example.org")
	}
	if config.Mail.SMTP.Address != "mail.example.com" {
		t.Errorf(
			"Mail.SMTP.Address == %v, want %v",
			config.Mail.SMTP.Address,
			"mail.example.com")
	}
	if config.Mail.SMTP.Password != "your_password" {
		t.Errorf(
			"Mail.SMTP.Password == %v, want %v",
			config.Mail.SMTP.Password,
			"your_password")
	}
	if config.Mail.SMTP.Port != 25 {
		t.Errorf("Mail.SMTP.Port == %v, want %v", config.Mail.SMTP.Port, 25)
	}
	if config.Mail.SMTP.UserName != "your.mail@example.org" {
		t.Errorf(
			"Mail.SMTP.UserName == %v, want %v",
			config.Mail.SMTP.UserName,
			"your.mail@example.org")
	}
	if config.Twitter.ConsumerKey != "your-consumer-key1" {
		t.Errorf(
			"Twitter.ConsumerKey == %v, want %v",
			config.Twitter.ConsumerKey,
			"your-consumer-key1")
	}
	if config.Twitter.ConsumerSecret != "your-consumer-secret1" {
		t.Errorf(
			"Twitter.ConsumerSecret == %v, want %v",
			config.Twitter.ConsumerSecret,
			"your-consumer-secret1")
	}
	if config.Twitter.AccessToken != "your-access-token1" {
		t.Errorf(
			"Twitter.AccessToken == %v, want %v",
			config.Twitter.AccessToken,
			"your-access-token1")
	}
	if config.Twitter.AccessTokenSecret != "your-access-token-secret1" {
		t.Errorf(
			"Twitter.AccessTokenSecret == %v, want %v",
			config.Twitter.AccessTokenSecret,
			"your-access-token-secret1")
	}
}
