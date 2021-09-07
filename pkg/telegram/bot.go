package telegram

import (
	"github.com/PavelStupnitski/telegram-bot-youtube/pkg/pocket"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient *pocket.Client
	redirectURL  string
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, redirectUrl string) *Bot {
	return &Bot{bot: bot, pocketClient: client, redirectURL: redirectUrl}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates, err := b.initUpdateChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) initUpdateChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
