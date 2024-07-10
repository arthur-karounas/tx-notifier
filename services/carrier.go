package services

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Service for sending notifications to recipients.
type Carrier struct {
	adminId string
}

func NewCarrier(adminId string) *Carrier {
	return &Carrier{
		adminId: adminId,
	}
}

// Send a specific message to users who are allowed to receive notifications.
func (c *Carrier) SendToAll(bot *tgbotapi.BotAPI, recipients []string, message string) {
	for _, recipient := range recipients {
		chatID, err := strconv.ParseInt(recipient, 10, 64)
		if err != nil {
			log.Printf("Carrier: Error parsing chat ID: %v", err)
			continue
		}

		SendMessage(bot, chatID, message)
	}

	log.Print("Carrier: Message sent to all recipients.")
}

// Send a specific message to the admin.
func (c *Carrier) SendToAdmin(bot *tgbotapi.BotAPI, message string) {
	adminChatID, err := strconv.ParseInt(c.adminId, 10, 64)
	if err != nil {
		log.Printf("Carrier: Error parsing admin ID: %v", err)
		return
	}

	SendMessage(bot, adminChatID, message)
}
