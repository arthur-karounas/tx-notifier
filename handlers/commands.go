package handlers

import (
	"txnotifier/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Routing incoming requests to the appropriate internal services.
func Commands(bot *tgbotapi.BotAPI, update tgbotapi.Update, fetcher *services.TronFetcher, carrier *services.Carrier, processor *services.TxProcessor, user *services.User, adminId string) {
	// Verify user as administrator.
	if !checkPermission(bot, update, adminId) {
		return
	}

	switch update.Message.Command() {
	case "help":
		sendHelpMessage(bot, update)

	case "add_user":
		user.AddUser(bot, update)
	case "delete_user":
		user.DeleteUser(bot, update)
	case "add_wallet":
		user.AddWallet(bot, update)
	case "delete_wallet":
		user.DeleteWallet(bot, update)

	case "start_notifications":
		processor.StartProcessing(bot)
	case "stop_notifications":
		processor.StopProcessing(bot, update)
	case "status":
		processor.CheckStatus(bot, update)
	}
}

// sendHelpMessage sends a help message with the list of available commands and their descriptions.
func sendHelpMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := `
<b>Available commands:</b>
- /help - Show this help message.

- /add_user [parameter] - Add a new user to the notification list.
- /delete_user [parameter] - Remove a user from the notification list.

- /add_wallet [parameter] - Add a new wallet to monitor.
- /delete_wallet [parameter] - Remove a wallet from monitoring.

- /start_notifications - Start transaction notifications.
- /stop_notifications - Stop transaction notifications.

- /status - Check the status of all monitored wallets.
`
	services.SendMessage(bot, update.Message.Chat.ID, message)
}
