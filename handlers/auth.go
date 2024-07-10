package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Chech if the user is an admin.
func checkPermission(bot *tgbotapi.BotAPI, update tgbotapi.Update, adminId string) bool {
	return adminId == fmt.Sprintf("%d", update.Message.Chat.ID)
}
