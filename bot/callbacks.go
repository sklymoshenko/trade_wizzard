package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCallbacks(update *tgbotapi.Update) {
	data := update.CallbackData()
	chatId := update.CallbackQuery.From.ID
	msg := tgbotapi.NewMessage(chatId, "Sorry error occuried.\nTry to start process from beginning /start")

	fmt.Println(data, msg)
}
