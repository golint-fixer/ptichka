# ptichka

Fetch tweets and sends by SMTP

[![Build Status](https://travis-ci.org/danil/ptichka.svg)](https://travis-ci.org/danil/ptichka)

## Description

`ptichka` should be run by cron periodically.

`ptichka` get your Twitter timeline messages and send them to you by email.

## Install

Go to http://apps.twitter.com and register application.

Then ...

## Limitations

### 200 messages

Load only last [200 messages][].
So if your have high flow then should process them often.

[200 messages]: https://dev.twitter.com/rest/reference/get/statuses/home_timeline#api-param-count

## Contributing

See the `CONTRIBUTING.md` file.

## License

See the `LICENSE` file.
