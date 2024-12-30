package bot

import (
	"log"
	"time"
	"trade_wizzard/ollama"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(update *tgbotapi.Update) {
	command := update.Message.Command()
	log.Println("Handling command", command)

	switch command {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! I am Trade Wizzard bot. I am providing a summarised information for each day . Use /help to see all available commands")
		b.sendMessage(msg)
	case "news":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Loading news for current day and making summary. Please wait...")
		b.sendMessage(msg)

		msg.Text = b.dayNewsSummary()
		b.sendMessage(msg)
		return
	case "stop":
	case "list":
	case "invite":
	}
}

func (b *Bot) dayNewsSummary() string {
	endTime := time.Now()
	startTime := time.Now().Add(-24 * time.Hour)
	newsResponse, err := b.ApiClient.GetNews(endTime, startTime)

	if err != nil || len(newsResponse.ReturnData) == 0 {
		log.Println("Error getting news", err)
		return "Sorry, I can't get news summary for you now. Try again later"
	}

	for _, news := range newsResponse.ReturnData {
		b.Ollama.SendMessage(ollama.OllamaMessage{Role: "user", Content: news.Title + "\n" + news.Body + "\n" + news.TimeString})
	}

	analysis := b.Ollama.Chat()

	return analysis.Content
}
