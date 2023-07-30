package main

import (
	"time"

	"github.com/MurashovVen/outsider-sdk/app/configuration"
)

type config struct {
	configuration.Default

	TelegramBotToken         string        `desc:"Auth token" split_words:"true"`
	TelegramBotUpdateTimeout time.Duration `desc:"Timeout for long polling" default:"0s" split_words:"true" `
}
