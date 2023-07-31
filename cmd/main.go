package main

import (
	"context"
	"errors"

	whether "github.com/MurashovVen/outsider-proto/whether/golang"
	"github.com/MurashovVen/outsider-sdk/app"
	"github.com/MurashovVen/outsider-sdk/app/configuration"
	"github.com/MurashovVen/outsider-sdk/app/logger"
	"github.com/MurashovVen/outsider-sdk/app/termination"
	"github.com/MurashovVen/outsider-sdk/grpc"
	tgclient "github.com/MurashovVen/outsider-sdk/tg"
	"go.uber.org/zap"

	"tg-bot/internal/tg"
)

func main() {
	cfg := new(config)
	configuration.MustProcessConfig(cfg)

	var (
		log = logger.MustCreateLogger(cfg.Env)

		telegramClient = tgclient.MustCreateAndConnect(cfg.TelegramBotToken)

		whetherClientConn = grpc.MustConnect(cfg.WhetherGRPCClientAddr, grpc.DefaultDialOptions(log)...)
	)

	application := app.New(
		log,
		app.AppendWorks(
			tg.New(
				telegramClient,
				whether.NewWhetherClient(whetherClientConn),
				tg.BotWithLogger(log),
			),
		),
	)

	if err := application.Run(context.Background()); err != nil && !errors.Is(err, termination.ErrStopped) {
		log.Error("running error: %s", zap.Error(err))
	}
}
