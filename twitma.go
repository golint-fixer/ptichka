package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	// "os"
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

func loadConfig(path string) (*twitmaConfig, error) {
	var config *twitmaConfig

	_, err := toml.DecodeFile(path, &config)

	return config, err
}

func loadIds(path string) ([]string, error) {
	var ids []string

	jsonBlob, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(jsonBlob, &ids)

	return ids, err
}

// func exists(path string) (bool, error) {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true, nil
// 	}
// 	if os.IsNotExist(err) {
// 		return false, nil
// 	}
// 	return true, err
// }

func main() {
	config, err := loadConfig(".twitmarc.toml")
	if err != nil {
		log.Fatal(err)
	}

	ids, err := loadIds(config.CacheFile)
	if err != nil {
		log.Fatal(err)
	}

	print(ids[0])
}
