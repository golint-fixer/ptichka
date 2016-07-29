package ptichka

import (
	"github.com/BurntSushi/toml"
)

type configs struct {
	Accounts []config
}

type config struct {
	CacheFile string `toml:"cache_file"`
	Label     string
	Verbose   bool
	Mail      struct {
		From struct {
			Address string
			Name    string
		}
		To struct {
			Address string
			Name    string
		}
		Method string
		SMTP   struct {
			Address  string
			Password string
			Port     int
			UserName string `toml:"user_name"`
		}
	}
	Twitter struct {
		ConsumerKey       string `toml:"consumer_key"`
		ConsumerSecret    string `toml:"consumer_secret"`
		AccessToken       string `toml:"access_token"`
		AccessTokenSecret string `toml:"access_token_secret"`
	}
}

func LoadConfig(path string) (*configs, error) {
	var configs *configs

	_, err := toml.DecodeFile(path, &configs)

	return configs, err
}
