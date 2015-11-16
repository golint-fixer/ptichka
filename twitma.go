package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type twitmaConfig struct {
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

var config *twitmaConfig

func loadConfig(path string) *twitmaConfig {
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func loadIds(jsonBlob []byte) []string {
	var ids []string
	if err := json.Unmarshal(jsonBlob, &ids); err != nil {
		log.Fatal(err)
	}
	return ids
}

func main() {
	config = loadConfig(".twitmarc.toml")
	file, err := ioutil.ReadFile(config.CacheFile)
	if err != nil {
		log.Fatal(err)
	}
	print(loadIds(file)[0])
}
