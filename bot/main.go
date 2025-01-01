package bot

import (
	"fmt"
	"log"
	"os"
	"trade_wizzard/api"
	"trade_wizzard/ollama"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/kejrak/xtb-client-go/xclient"
)

type EventContext struct {
	Status         string
	SelectedUserId int64
}

type ChatContext struct {
	ID           int64
	EventContext *EventContext
	User         *tgbotapi.User
}

type Bot struct {
	ApiToken           string
	ApiClient          *xclient.Client
	Ollama             *ollama.Ollama
	CurrentChatContext map[int64]*ChatContext
	Url                string
}

var tgBot *tgbotapi.BotAPI

func (b *Bot) LoadApiKey() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Cant load env file")
	}

	apiKey := os.Getenv("T_API_TOKEN")

	b.ApiToken = apiKey
}

func (b *Bot) Start() error {
	b.LoadApiKey()

	client, err := api.XtbClient()

	if err != nil {
		log.Println("Cant create xtb client")
		return err
	}

	b.ApiClient = client
	b.Ollama = ollama.NewOllama()

	tgBot, err = tgbotapi.NewBotAPI(b.ApiToken)

	if err != nil {
		log.Fatalf("Cant create tgbot: %v", err)
	}

	log.Printf("Authorized on account %s", tgBot.Self.UserName)
	b.Url = fmt.Sprintf("https://t.me/%s", tgBot.Self.UserName)

	updateConfig := tgbotapi.UpdateConfig{Offset: 0, Timeout: 60}
	updates := tgBot.GetUpdatesChan(updateConfig)

	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.CallbackQuery != nil {
			go b.handleCallbacks(&update)
		} else if update.Message.IsCommand() {
			go b.handleCommands(&update)
		} else {
			go b.handleTextReplies(&update)
		}
	}
}

func (b *Bot) sendMessage(msg tgbotapi.Chattable) {
	if _, err := tgBot.Send(msg); err != nil {
		log.Panicf("Error sending message %v", err)
	}
}
