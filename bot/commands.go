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

		newsChan := make(chan string)
		go b.dayNewsSummary(newsChan)

		for news := range newsChan {
			msg.Text = news
			b.sendMessage(msg)
		}

		msg.Text = "News summary completed"
		b.sendMessage(msg)
		return
	case "stop":
	case "list":
	case "invite":
	}
}

func (b *Bot) dayNewsSummary(newsChan chan<- string) {
	defer close(newsChan)

	endTime := time.Now()
	startTime := time.Now().Add(-48 * time.Hour)
	newsResponse, err := b.ApiClient.GetNews(endTime, startTime)

	if err != nil {
		log.Println("Error getting news", err)
		newsChan <- "Sorry, I can't get news summary for you now. Try again later"
		return
	}

	if len(newsResponse.ReturnData) == 0 {
		newsChan <- "No news for today"
		return
	}

	prompt := ""
	b.Ollama.SendMessage(ollama.OllamaMessage{Role: "system", Content: "Provide brief summary for each news article. Make sure to include the title of the article in the summary."})
	for _, news := range newsResponse.ReturnData {
		prompt += "Title:" + news.Title + "\n" + "Body: " + news.Body + "\n" + "-------" + "\n"
		b.Ollama.SendMessage(ollama.OllamaMessage{Role: "user", Content: prompt})
		analysis := b.Ollama.Chat()
		newsChan <- analysis.Content
		// newsChan <- news.Title
	}
}
