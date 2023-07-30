package tg

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *Bot) Start(ctx context.Context) error {
	b.logger.Debug("connecting")

	if err := b.connect(); err != nil {
		return err
	}

	updChan := b.tg.GetUpdatesChan(
		tgbotapi.UpdateConfig{
			Offset:  0, // todo need to store processed offsets ?
			Limit:   0,
			Timeout: int(b.updateTimeout.Seconds()),
		},
	)

	b.logger.Debug("start receiving messages")

	for {
		select {
		case <-ctx.Done():
			b.logger.Debug("stopping")

			b.tg.StopReceivingUpdates()

			b.logger.Debug("stopped")

			return nil

		case upd := <-updChan:
			if upd.Message != nil {
				b.logger.Debug("got message", zap.String("msg", upd.Message.Text))
			}
		}
	}
}

func (b *Bot) connect() error {
	tgBot, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return fmt.Errorf("подключение к апи: %w", err)
	}

	b.tg = tgBot

	return nil
}
