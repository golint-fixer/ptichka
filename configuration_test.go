package ptichka

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestConfigurations(t *testing.T) {
	var configs *Configurations
	_, err := toml.DecodeFile(".ptichkarc.toml.example", &configs)
	if err != nil {
		t.Error(err)
	}

	config := configs.Accounts[0]

	var wont string
	var wontInt int

	wont = "/path/to/cache/file1.json"
	if config.CacheFile != wont {
		t.Errorf("CacheFile == %v, want %v", config.CacheFile, wont)
	}

	wont = "[twitter account 1] "
	if config.Label != wont {
		t.Errorf("Label == %v, want %v", config.Label, wont)
	}

	wont = "noreply@example.com"
	if config.Mail.From.Address != wont {
		t.Errorf("Mail.From.Address == %v, want %v", config.Mail.From.Address, wont)
	}

	wont = "your.mail+twitter1@example.org"
	if config.Mail.To.Address != wont {
		t.Errorf("Mail.To.Address == %v, want %v", config.Mail.To.Address, wont)
	}

	wont = "mail.example.com"
	if config.Mail.SMTP.Address != wont {
		t.Errorf("Mail.SMTP.Address == %v, want %v", config.Mail.SMTP.Address, wont)
	}

	wont = "your_password"
	if config.Mail.SMTP.Password != wont {
		t.Errorf("Mail.SMTP.Password == %v, want %v", config.Mail.SMTP.Password, wont)
	}

	wontInt = 25
	if config.Mail.SMTP.Port != wontInt {
		t.Errorf("Mail.SMTP.Port == %v, want %v", config.Mail.SMTP.Port, wontInt)
	}

	wont = "your.mail@example.org"
	if config.Mail.SMTP.UserName != wont {
		t.Errorf(
			"Mail.SMTP.UserName == %v, want %v",
			config.Mail.SMTP.UserName,
			wont)
	}

	wont = "your-consumer-key1"
	if config.Twitter.ConsumerKey != wont {
		t.Errorf(
			"Twitter.ConsumerKey == %v, want %v",
			config.Twitter.ConsumerKey,
			wont)
	}

	wont = "your-consumer-secret1"
	if config.Twitter.ConsumerSecret != wont {
		t.Errorf(
			"Twitter.ConsumerSecret == %v, want %v",
			config.Twitter.ConsumerSecret,
			wont)
	}

	wont = "your-access-token1"
	if config.Twitter.AccessToken != wont {
		t.Errorf(
			"Twitter.AccessToken == %v, want %v",
			config.Twitter.AccessToken,
			wont)
	}

	wont = "your-access-token-secret1"
	if config.Twitter.AccessTokenSecret != wont {
		t.Errorf(
			"Twitter.AccessTokenSecret == %v, want %v",
			config.Twitter.AccessTokenSecret,
			wont)
	}
}
