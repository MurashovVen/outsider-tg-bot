package tg

import (
	"context"
	"time"

	"github.com/MurashovVen/outsider-sdk/app"
	"github.com/MurashovVen/outsider-sdk/app/logger"
	"github.com/MurashovVen/outsider-sdk/tg"

	whether "github.com/MurashovVen/outsider-proto/whether/golang"
)

type (
	Bot struct {
		tg *tg.Client

		whetherService whether.WhetherClient

		logger *logger.Logger

		updateTimeout time.Duration
	}

	BotOption func(*Bot)
)

func New(tg *tg.Client, whetherService whether.WhetherClient, options ...BotOption) *Bot {
	b := &Bot{
		tg:             tg,
		whetherService: whetherService,
		logger:         logger.NewNop(),
	}

	for _, opt := range options {
		opt(b)
	}

	return b
}

func BotWithLogger(logger *logger.Logger) BotOption {
	return func(bot *Bot) {
		bot.logger = logger.Named("TelegramBot")
	}
}

var (
	_ app.Work = (*Bot)(nil)
)

func (b *Bot) Runner(ctx context.Context) func() error {
	return func() error {
		return b.start(ctx)
	}
}

func (b *Bot) Name() string {
	return "TelegramBot"
}
