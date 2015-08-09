require 'mail'
require 'net/http'
require 'pathname'
require 'rubygems'
require 'twitter'
require 'yaml'

ROOT = Pathname(__FILE__).dirname
CONFIG = ROOT.join('.thtmrc.yml')
READED_IDS = ROOT.join('.thtm.yml')

config = File.exist?(CONFIG) ? YAML.load_file(CONFIG) : {}

CHARSET = config['charset'] || 'UTF-8'
COUNT = config['count'] || 200
LABEL = config['label'] || 'twitter'
VERBOSE = config.key?('verbose') ? config['verbose'] : false

mailrc = config['mail'] || {}
FROM = mailrc['from'] || 'noreply@example.org'
TO = mailrc['to'] || fail('You should provide your email adress')

twitterrc = config['twitter'] || {}
CONSUMER_KEY = twitterrc['consumer_key']
CONSUMER_SECRET = twitterrc['consumer_secret']
ACCESS_TOKEN = twitterrc['access_token']
ACCESS_TOKEN_SECRET = twitterrc['access_token_secret']

readed_ids = File.exist?(READED_IDS) ? YAML.load_file(READED_IDS) : []

puts("Old messages: #{readed_ids.size}") if VERBOSE

client = Twitter::REST::Client.new do |c|
  c.consumer_key = CONSUMER_KEY
  c.consumer_secret = CONSUMER_SECRET
  c.access_token = ACCESS_TOKEN
  c.access_token_secret = ACCESS_TOKEN_SECRET
end

messages = client.home_timeline(count: COUNT)

messages.sort_by!(&:created_at)

puts("Loaded messages: #{messages.size}") if VERBOSE

if mailrc['method']
  mail_options = {}
  mail_options.merge(mailrc['smtp']) if mailrc['method'] == 'smtp'

  Mail.defaults do
    delivery_method mailrc['method'], mail_options
  end
end

messages.each do |message|
  next if readed_ids.include?(message.id)

  s = format('@%{screen_name} %{created_at}',
             screen_name: message.user.screen_name,
             created_at: message.created_at)

  b = format("%{text}\n--\nhttps://twitter.com/%{screen_name}/status/%{id}",
             text: message.text,
             id: message.id,
             screen_name: message.user.screen_name)

  mail = Mail.new do
    from FROM
    to TO
    subject "[#{LABEL}] #{s}"

    if message.media.any?
      message.media.each do |media|
        url = URI.parse(media.media_url)
        attachments[Pathname.new(url.to_s).basename.to_s] = Net::HTTP.get(url)
      end

      text_part do
        content_type "text/plain; charset=#{CHARSET.downcase}"
        body b
      end
    else
      content_type "text/plain; charset=#{CHARSET.downcase}"
      body b
    end
  end

  puts("Sending: #{s}") if VERBOSE
  mail.deliver!
end

readed_ids.push(*messages.map(&:id)).uniq!

puts("Total messages: #{readed_ids.size}") if VERBOSE

File.open(READED_IDS, 'w') { |f| f.write readed_ids.to_yaml }
