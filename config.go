package main

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	CacheFile string `toml:"cache_file"`
	Count     int
	Label     string
	Verbose   bool
	Mail      struct {
		From   string
		To     string
		Method string
		SMTP   struct {
			Address        string
			Authentication string
			Password       string
			Port           int
			SSL            bool
			TLS            bool
			UserName       string `toml:"user_name"`
		}
	}
	Twitter struct {
		ConsumerKey       string `toml:"consumer_key"`
		ConsumerSecret    string `toml:"consumer_secret"`
		AccessToken       string `toml:"access_token"`
		AccessTokenSecret string `toml:"access_token_secret"`
	}
}

func loadConfig(path string) (*config, error) {
	var config *config

	_, err := toml.DecodeFile(path, &config)

	return config, err
}
