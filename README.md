# twitma

Twitter timeline mailer.

## Description

This utility should run periodically by cron.

Get your timeline messages
and send them to you by mail.

## Install

Go to http://apps.twitter.com and register application.
Then

```sh
cd path/to/twitma
bundle install
cp --interactive .twitmarc.yml.example .twitmarc.yml
```

Then set mail and twitter api configuration in `.twitma.yml`.

Then create cron task like that:
`bundle exec 'ruby path/to/twitma.rb'`

## Limitations/Bugs

### 200 messages

Load only last 200 messages.
So if your have high flow then should process them often.

### Sort order

Sent mails can come in the wrong order.

## License

The MIT License (MIT)

Copyright (c) 2015 Danil Kutkevich <danil@kutkevich.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
