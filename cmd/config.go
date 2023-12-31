package main

import (
	"time"

	"github.com/MurashovVen/outsider-sdk/app/configuration"
)

type config struct {
	configuration.Default
	configuration.TelegramClient

	WhetherGRPCClientAddr    string        `desc:"Address of whether service" default:"whether:5000" split_words:"true"`
	TelegramBotUpdateTimeout time.Duration `desc:"Timeout for long polling" default:"0s" split_words:"true"`
}
