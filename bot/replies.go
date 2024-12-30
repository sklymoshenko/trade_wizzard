package bot

import (
	"fmt"
	"log"
	"trade_wizzard/ollama"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleTextReplies(update *tgbotapi.Update) {
	log.Println("Handling reply")
	fmt.Println(update.Message.Text)

	b.Ollama.SendMessage(ollama.OllamaMessage{Role: "user", Content: update.Message.Text})
	response := b.Ollama.Chat()
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response.Content)
	b.sendMessage(msg)
}
