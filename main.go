package main

import (
	"trade_wizzard/bot"
)

func main() {
	tgBot := &bot.Bot{}
	tgBot.Start()

	// api.StartApiServer()
}
