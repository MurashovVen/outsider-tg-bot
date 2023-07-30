package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/MurashovVen/outsider-sdk/app"
	"github.com/MurashovVen/outsider-sdk/app/logger"
	"github.com/MurashovVen/outsider-sdk/app/termination"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"tg-bot/internal/tg"
)

func main() {
	cfg := new(config)
	app.MustProcessConfig(cfg)

	log := logger.MustCreateLogger(cfg.Env)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(botRunner(ctx, cfg, log))

	eg.Go(
		termination.NewWaiter(
			termination.WaiterWithLogger(log),
		).WaitFunc(ctx),
	)

	eg.Go(log.SyncWaiter(ctx))

	if err := eg.Wait(); err != nil && !errors.Is(err, termination.ErrStopped) {
		log.Error("running error: %s", zap.Error(err))
	}
}

func botRunner(ctx context.Context, cfg *config, log *logger.Logger) func() error {
	return func() error {
		b := tg.New(
			cfg.TelegramBotToken,
			cfg.TelegramBotUpdateTimeout,
			tg.BotWithLogger(log),
		)

		if err := b.Start(ctx); err != nil {
			return fmt.Errorf("starting bot: %w", err)
		}

		return nil
	}
}
