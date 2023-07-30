package main

import (
	"time"

	"tg-bot/pkg/app"
)

type config struct {
	Env                      app.Environment `desc:"(development)" default:"development" split_words:"true"`
	TelegramBotToken         string          `desc:"Auth token" split_words:"true"`
	TelegramBotUpdateTimeout time.Duration   `desc:"Timeout for long polling" default:"0s" split_words:"true" `
}
