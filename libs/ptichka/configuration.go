package ptichka

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

// Configurations ia an structure with collection of timelines configurations.
type Configurations struct {
	Accounts []*configuration
}

type configuration struct {
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

func (config *configuration) mailFrom() string {
	from := mail.Address{
		Name:    config.Mail.From.Name,
		Address: config.Mail.From.Address}

	return from.String()
}

func (config *configuration) mailTo() string {
	to := mail.Address{
		Name:    config.Mail.To.Name,
		Address: config.Mail.To.Address}

	return to.String()
}

// Sender is an interface with send method.
type Sender interface {
	send(config *configuration, msg []byte) error
}

// SMTPSender has a method to send emails
// realized through the net.smtp package.
type SMTPSender struct{}

// send method sends messages through SMTP.
func (s *SMTPSender) send(config *configuration, msg []byte) error {
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", config.Mail.SMTP.Address, config.Mail.SMTP.Port),
		smtp.PlainAuth(
			"",
			config.Mail.SMTP.UserName,
			config.Mail.SMTP.Password,
			config.Mail.SMTP.Address),
		config.mailFrom(),
		[]string{config.mailTo()},
		msg)
}
