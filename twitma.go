package main

import (
	"log"

	"github.com/BurntSushi/toml"
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

func main() {
	config = loadConfig(".twitmarc.toml")

	print(config.Mail.SMTP.Address)
}

// require 'mail'
// require 'net/http'
// require 'pathname'
// require 'twitter'

// ROOT = Pathname(__FILE__).dirname
// CONFIG = ROOT.join('.twitmarc.yml')
// READED_IDS = ROOT.join('.twitma.yml')

// config = File.exist?(CONFIG) ? YAML.load_file(CONFIG) : {}

// CHARSET = config['charset'] || 'UTF-8'
// COUNT = config['count'] || 200
// LABEL = config['label'] || 'twitter'
// VERBOSE = config.key?('verbose') ? config['verbose'] : false

// mailrc = config['mail'] || {}
// FROM = mailrc['from'] || 'noreply@example.org'
// TO = mailrc['to'] || fail('You should provide your email adress')

// twitterrc = config['twitter'] || {}
// CONSUMER_KEY = twitterrc['consumer_key']
// CONSUMER_SECRET = twitterrc['consumer_secret']
// ACCESS_TOKEN = twitterrc['access_token']
// ACCESS_TOKEN_SECRET = twitterrc['access_token_secret']

// readed_ids = File.exist?(READED_IDS) ? YAML.load_file(READED_IDS) : []

// puts("Old messages: #{readed_ids.size}") if VERBOSE

// client = Twitter::REST::Client.new do |c|
//   c.consumer_key = CONSUMER_KEY
//   c.consumer_secret = CONSUMER_SECRET
//   c.access_token = ACCESS_TOKEN
//   c.access_token_secret = ACCESS_TOKEN_SECRET
// end

// messages = client.home_timeline(count: COUNT)

// # File.open('/tmp/foo_bar.yml', 'w') { |f| f.write messages.to_yaml }
// # messages = YAML.load_file('/tmp/foo_bar.yml')

// messages.sort_by!(&:created_at)

// puts("Loaded messages: #{messages.size}") if VERBOSE

// if mailrc['method']
//   mail_options = {}
//   mail_options.merge(mailrc['smtp']) if mailrc['method'] == 'smtp'

//   Mail.defaults do
//     delivery_method mailrc['method'], mail_options
//   end
// end

// messages.each do |message|
//   next if readed_ids.include?(message.id)

//   s = format('@%{name} %{date}',
//              name: message.user.screen_name,
//              date: message.created_at)

//   b = format("@%{name}\n\n%{text}\n\n%{url}",
//              name: message.user.screen_name,
//              text: message.text,
//              screen_name: message.user.screen_name,
//              url: format('https://twitter.com/%{name}/status/%{id}',
//                          name: message.user.screen_name,
//                          id: message.id))

//   mail = Mail.new do
//     from FROM
//     to TO
//     subject "[#{LABEL}] #{s}"

//     if message.media.any?
//       message.media.each do |media|
//         url = URI.parse(media.media_url)
//         attachments[Pathname.new(url.to_s).basename.to_s] = Net::HTTP.get(url)
//       end

//       text_part do
//         content_type "text/plain; charset=#{CHARSET.downcase}"
//         body b
//       end
//     else
//       content_type "text/plain; charset=#{CHARSET.downcase}"
//       body b
//     end
//   end

//   puts("Sending: #{s}") if VERBOSE
//   mail.deliver!
// end

// readed_ids.push(*messages.map(&:id)).uniq!

// puts("Total messages: #{readed_ids.size}") if VERBOSE

// File.open(READED_IDS, 'w') { |f| f.write readed_ids.to_yaml }
