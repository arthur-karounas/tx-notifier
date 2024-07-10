package main

import (
	"log"
	"os"
	"txnotifier/clients"
	"txnotifier/handlers"
	"txnotifier/services"

	"github.com/joho/godotenv"
)

func main() {
	// Loading the .env configuration.
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Config: Error loading .env file.")
	}

	// Initialization of the services.
	carrier := services.NewCarrier(os.Getenv("ADMIN_CHAT_ID"))
	user := services.NewUser(carrier, os.Getenv("USER_JSON"))
	fetcher := services.NewTronFetcher(carrier, os.Getenv("TRON_API_ENDPOINT"))
	processor := services.NewTxProcessor(fetcher, carrier, user)

	// Initialization of the telegram bot.
	bot := clients.Init(os.Getenv("TG_BOT_TOKEN"))

	// Start processing incoming requests.
	handlers.Init(bot, fetcher, carrier, processor, user, os.Getenv("ADMIN_CHAT_ID"))
}
