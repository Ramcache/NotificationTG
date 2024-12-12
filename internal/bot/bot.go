package bot

import (
	"NotificationTG/internal/db"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	TelegramBot *tgbotapi.BotAPI
	DB          *db.DB
}

func NewBot(token string, database *db.DB) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		TelegramBot: bot,
		DB:          database,
	}, nil
}

func (b *Bot) Run() {
	b.TelegramBot.Debug = true
	log.Printf("Authorized on account %s", b.TelegramBot.Self.UserName)

	updates := b.TelegramBot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message != nil {
			b.HandleMessage(update.Message)
		}
	}
}
