package tg

import (
	"context"

	whether "github.com/MurashovVen/outsider-proto/whether/golang"
	"github.com/MurashovVen/outsider-sdk/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *Bot) start(ctx context.Context) error {
	b.logger.Debug("connecting")

	updChan := b.tg.GetUpdatesChan(
		tgbotapi.UpdateConfig{
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
			if err := b.processUpdate(ctx, upd); err != nil {
				b.logger.Error("processing update", zap.Error(err))
			}
		}
	}
}

func (b *Bot) processUpdate(ctx context.Context, upd tgbotapi.Update) error {
	switch {
	case upd.Message != nil:
		return b.processMessage(upd.Message)

	case upd.CallbackQuery != nil:
		return b.processCallbackQuery(ctx, upd.CallbackQuery)

	default:
		return ErrUnsupportedUpdate
	}
}

func (b *Bot) processCallbackQuery(ctx context.Context, cd *tgbotapi.CallbackQuery) error {
	action := entities.ActionTypeParseString(cd.Data)

	switch {
	case action.IsWhetherType():
		_, err := b.whetherService.ActionProcess(ctx,
			&whether.ActionProcessRequest{
				FromChatId: cd.Message.Chat.ID,
				Action:     cd.Data,
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

// todo вынести в сервис тоже
func (b *Bot) sendStart(msg *tgbotapi.Message) error {
	cbData := entities.ActionWhetherConfigure

	_, err := b.tg.Send(
		&tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: msg.Chat.ID,
				ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
						{
							{
								Text:         "Сконфигурировать погоду",
								CallbackData: &cbData,
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
