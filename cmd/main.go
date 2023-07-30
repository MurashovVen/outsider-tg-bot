package main

import (
	"context"
	"errors"

	"github.com/MurashovVen/outsider-sdk/app"
	"github.com/MurashovVen/outsider-sdk/app/configuration"
	"github.com/MurashovVen/outsider-sdk/app/logger"
	"github.com/MurashovVen/outsider-sdk/app/termination"
	"go.uber.org/zap"

	"tg-bot/internal/tg"
)

func main() {
	cfg := new(config)
	configuration.MustProcessConfig(cfg)

	log := logger.MustCreateLogger(cfg.Env)

	application := app.New(
		log,
		app.AppendWorks(
			tg.New(
				cfg.TelegramBotToken,
				cfg.TelegramBotUpdateTimeout,
				tg.BotWithLogger(log),
			),
		),
	)

	if err := application.Run(context.Background()); err != nil && !errors.Is(err, termination.ErrStopped) {
		log.Error("running error: %s", zap.Error(err))
	}
}
