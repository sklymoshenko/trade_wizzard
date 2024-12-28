package ollama

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
)

type OllamaMessage = api.Message

type Ollama struct {
	Client   *api.Client
	Messages []OllamaMessage
}

func NewOllama() *Ollama {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []OllamaMessage{
		{
			Role:    "system",
			Content: "You are a very proffesional and smart finance specialist. Provide brief analysis of the current market situation depending on news and events. It need to be understandable for a beginner. Always put at the end the news source of what lead you to some conclusion. For simple questions please answer without tradind analysis.",
		},
	}

	return &Ollama{Client: client, Messages: messages}
}

func (o *Ollama) SendMessage(msg OllamaMessage) {
	o.Messages = append(o.Messages, msg)
}

func (o *Ollama) Chat() OllamaMessage {
	responseChan := make(chan OllamaMessage)
	stream := false
	defer close(responseChan)

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "llama3.2",
		Messages: o.Messages,
		Stream:   &stream,
	}

	respFunc := func(resp api.ChatResponse) error {
		o.Messages = append(o.Messages, resp.Message)
		responseChan <- resp.Message
		return nil
	}

	go func() {
		err := o.Client.Chat(ctx, req, respFunc)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return <-responseChan
}

func (o *Ollama) StreamChat(responseChan chan<- OllamaMessage) {
	var fullMessage OllamaMessage

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "llama3.2",
		Messages: o.Messages,
	}

	respFunc := func(resp api.ChatResponse) error {
		fullMessage.Content += resp.Message.Content
		responseChan <- resp.Message
		return nil
	}

	go func() {
		err := o.Client.Chat(ctx, req, respFunc)
		if err != nil {
			log.Fatal(err)
		}
		o.Messages = append(o.Messages, fullMessage)
		close(responseChan)
	}()
}
