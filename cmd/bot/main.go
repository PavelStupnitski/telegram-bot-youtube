package main

import (
	"github.com/PavelStupnitski/telegram-bot-youtube/pkg/config"
	"github.com/PavelStupnitski/telegram-bot-youtube/pkg/pocket"
	"github.com/PavelStupnitski/telegram-bot-youtube/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	uniqueKeys, err := config.ReadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(uniqueKeys.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(uniqueKeys.ConsumerKey)
	if err != nil {
		log.Fatal(err)
	}
	telegramBot := telegram.NewBot(bot, pocketClient, "http://localhost")
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
