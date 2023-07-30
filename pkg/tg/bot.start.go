package tg

import (
	"context"
	"fmt"
	"strings"

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
			Offset:  0,
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
			if err := b.processUpdate(upd); err != nil {
				b.logger.Error("processing update", zap.Error(err))
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

func (b *Bot) processUpdate(upd tgbotapi.Update) error {
	switch {
	case upd.Message != nil:
		return b.processMessage(upd.Message)

	case upd.CallbackQuery != nil:
		return b.processCallbackQuery(upd.CallbackQuery)

	default:
		return ErrUnsupportedUpdate
	}
}

func (b *Bot) processCallbackQuery(cd *tgbotapi.CallbackQuery) error {
	switch {
	case cd.Data == CallbackDataConfigureWhether:
		buttons, err := callbackDataConfigureWhetherTemperatureCreateButtons(-40, 40)
		if err != nil {
			return fmt.Errorf("creating buttons callback data configure whether buttons: %w", err)
		}

		_, err = b.tg.Send(
			&tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: cd.Message.Chat.ID,
					ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
						InlineKeyboard: buttons,
					},
				},
				Text: `Выберете критическое значение температуры`,
			},
		)
		return err

	case strings.Contains(cd.Data, CallbackDataConfigureWhetherTemperature):
		_, err := b.tg.Send(
			&tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: cd.Message.Chat.ID,
				},
				Text: `Вы подписались на обновления и сконфигурировали критическую температуру. Спасибо)`,
			},
		)
		return err

	default:
		return ErrUnsupportedCallbackData
	}
}

func (b *Bot) processMessage(msg *tgbotapi.Message) error {
	switch {
	case msg.IsCommand():
		return b.processCommand(msg)

	default:
		return ErrUnsupportedMessage
	}
}

func (b *Bot) processCommand(msg *tgbotapi.Message) error {
	if !msg.IsCommand() {
		return ErrMsgIsNotCommand
	}

	switch msg.Command() {
	case "start":
		return b.sendStart(msg)

	case "help":
		fallthrough
	default:
		return b.sendHelp(msg)
	}
}

func (b *Bot) sendStart(msg *tgbotapi.Message) error {
	_, err := b.tg.Send(
		&tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: msg.Chat.ID,
				ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
						{
							{
								Text:         "Сконфигурировать погоду",
								CallbackData: &CallbackDataConfigureWhether,
							},
						},
					},
				},
			},
			Text: `Что вы хотите сконфигурировать?`,
		},
	)
	return err
}

func (b *Bot) sendHelp(msg *tgbotapi.Message) error {
	_, err := b.tg.Send(&tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: msg.Chat.ID,
		},
		Text: `Список доступных комманд:
/start - начать взаимодействие с ботом`,
	})
	return err
}
