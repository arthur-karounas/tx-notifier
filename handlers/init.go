package handlers

import (
	"txnotifier/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start processing incoming requests to the bot.
func Init(bot *tgbotapi.BotAPI, fetcher *services.TronFetcher, carrier *services.Carrier, processor *services.TxProcessor, user *services.User, adminId string) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Start processing of incoming commands.
	for update := range updates {
		if update.Message.IsCommand() {
			Commands(bot, update, fetcher, carrier, processor, user, adminId)
		}
	}
}
