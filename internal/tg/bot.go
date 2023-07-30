package tg

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MurashovVen/outsider-sdk/app/logger"
)

type (
	Bot struct {
		tg *tgbotapi.BotAPI

		logger *logger.Logger

		token         string
		updateTimeout time.Duration
	}

	BotOption func(*Bot)
)

func New(token string, updateTimeout time.Duration, options ...BotOption) *Bot {
	b := &Bot{
		token:         token,
		updateTimeout: updateTimeout,
		logger:        logger.NewNop(),
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
