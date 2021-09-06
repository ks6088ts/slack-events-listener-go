[![test](https://github.com/ks6088ts/slack-events-listener-go/workflows/test/badge.svg)](https://github.com/ks6088ts/slack-events-listener-go/actions/workflows/test.yml)
[![release](https://github.com/ks6088ts/slack-events-listener-go/workflows/release/badge.svg)](https://github.com/ks6088ts/slack-events-listener-go/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ks6088ts/slack-events-listener-go)](https://goreportcard.com/report/github.com/ks6088ts/slack-events-listener-go)
[![GoDoc](https://godoc.org/github.com/ks6088ts/slack-events-listener-go?status.svg)](https://godoc.org/github.com/ks6088ts/slack-events-listener-go)

# slack-events-listener-go

## Usage

```bash
A web server which handles events received from Slack Events API written in Go

Usage:
  cli [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run a slack event listener
  version     Print version info

Flags:
      --config string   config file (default is $HOME/.cli.yaml)
  -h, --help            help for cli
  -t, --toggle          Help message for toggle

Use "cli [command] --help" for more information about a command.
```

## Get Started

1. Setup Slack app settings to subscribe event types you need.
2. Setup Google Sheets API and service accounts
3. Run server with the following commands
4. Setup a proxy server and Set Request URL.
   For example, if you use `ngrok`, you should set `https://YOURDOMAIN.ngrok.io/slack/events` to Request URL in Slack apps settings.

```
cli run \
    --secret SLACK_SIGNING_SECRET \
    --token SLACK_BOT_TOKEN \
    --credentials PATH_TO_CREDENTIALS_FOR_GOOGLE \
    --sheetId SHEET_ID
```

## Reference

- [Using the Slack Events API](https://api.slack.com/apis/connections/events-api)
- [Service accounts](https://cloud.google.com/iam/docs/service-accounts)
