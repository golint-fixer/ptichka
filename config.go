package ptichka

import (
	"net/mail"

	"github.com/BurntSushi/toml"
)

// Configs ia an structure with collection of timelines configs.
type Configs struct {
	Accounts []config
}

type config struct {
	CacheFile string `toml:"cache_file"`
	Label     string
	Verbose   bool
	LogFile   string `toml:"log_file"`
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

func (config config) mailFrom() string {
	from := mail.Address{
		Name:    config.Mail.From.Name,
		Address: config.Mail.From.Address}

	return from.String()
}

func (config config) mailTo() string {
	to := mail.Address{
		Name:    config.Mail.To.Name,
		Address: config.Mail.To.Address}

	return to.String()
}

// LoadConfig load TOML config files.
func LoadConfig(path string) (*Configs, error) {
	var configs *Configs

	_, err := toml.DecodeFile(path, &configs)

	return configs, err
}
