package tg

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	CallbackDataConfigureWhether            = "configure-whether"
	CallbackDataConfigureWhetherTemperature = "configure-whether-temperature"
)

func callbackDataConfigureWhetherTemperatureCreateButtons() [][]tgbotapi.InlineKeyboardButton {
	var (
		rowLen = 5

		from, to = -40, 40

		res = make([][]tgbotapi.InlineKeyboardButton, 0, (to-from+1)/rowLen+1)

		currRow = make([]tgbotapi.InlineKeyboardButton, 0, rowLen)
	)

	for ; from <= to; from++ {
		cbData := CallbackDataConfigureWhetherTemperature + ":" + strconv.Itoa(from)
		currRow = append(currRow,
			tgbotapi.InlineKeyboardButton{
				Text:         strconv.Itoa(from),
				CallbackData: &cbData,
			},
		)

		if len(currRow) == rowLen {
			res = append(res, currRow)

			currRow = make([]tgbotapi.InlineKeyboardButton, 0, rowLen)
		}
	}

	return res
}
