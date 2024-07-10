package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// User represents a user receiving notifications.
type User struct {
	carrier         *Carrier // Carrier for sending messages.
	filename        string   // Filename for saving user data.
	Recipients      []string // List of recipients to notify.
	WalletAddresses []string // List of wallet addresses to monitor.
}

// NewUser initializes a new user with the given carrier and filename.
func NewUser(carrier *Carrier, filename string) *User {
	user := &User{
		carrier:         carrier,
		filename:        filename,
		Recipients:      []string{},
		WalletAddresses: []string{},
	}

	// Load user data from file.
	err := user.loadFromFile()
	if err != nil {
		log.Printf("User: Error loading user data: %v. Using default values.", err)
	}

	log.Println("User: Using the saved user configuration.")

	return user
}

// AddUser adds a user to receive notifications.
func (u *User) AddUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	chatID := strings.TrimSpace(update.Message.CommandArguments())
	if chatID == "" {
		message := "Please specify a parameter after command."
		u.carrier.SendToAdmin(bot, message)
		return nil
	}

	u.Recipients = append(u.Recipients, chatID)

	message := "User was successfully added."
	log.Print("User: ", message)
	u.carrier.SendToAdmin(bot, message)

	if err := u.saveToFile(); err != nil {
		log.Printf("Error saving user data: %v", err)
	}

	return nil
}

// DeleteUser removes a user from receiving notifications.
func (u *User) DeleteUser(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	chatID := strings.TrimSpace(update.Message.CommandArguments())
	if chatID == "" {
		message := "Please specify a parameter after command."
		u.carrier.SendToAdmin(bot, message)
		return nil
	}

	u.Recipients = removeString(u.Recipients, chatID)

	message := "User was successfully deleted."
	log.Print("User: ", message)
	u.carrier.SendToAdmin(bot, message)

	if err := u.saveToFile(); err != nil {
		log.Printf("Error saving user data: %v", err)
	}

	return nil
}

// AddWallet adds a wallet to monitor.
func (u *User) AddWallet(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	walletAddress := strings.TrimSpace(update.Message.CommandArguments())
	if walletAddress == "" {
		message := "Please specify a parameter after command."
		u.carrier.SendToAdmin(bot, message)
		return nil
	}

	u.WalletAddresses = append(u.WalletAddresses, walletAddress)

	message := "Wallet was successfully added."
	log.Print("User: ", message)
	u.carrier.SendToAdmin(bot, message)

	if err := u.saveToFile(); err != nil {
		log.Printf("Error saving user data: %v", err)
	}

	return nil
}

// DeleteWallet removes a monitored wallet.
func (u *User) DeleteWallet(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	walletAddress := strings.TrimSpace(update.Message.CommandArguments())
	if walletAddress == "" {
		message := "Please specify a parameter after command."
		u.carrier.SendToAdmin(bot, message)
		return nil
	}

	u.WalletAddresses = removeString(u.WalletAddresses, walletAddress)

	message := "Wallet was successfully deleted."
	log.Print("User: ", message)
	u.carrier.SendToAdmin(bot, message)

	if err := u.saveToFile(); err != nil {
		log.Printf("Error saving user data: %v", err)
	}

	return nil
}

// saveToFile saves the user data to a JSON file.
func (u *User) saveToFile() error {
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(u.filename, data, 0644)
}

// LoadFromFile loads the user data from a JSON file.
func (u *User) loadFromFile() error {
	if _, err := os.Stat(u.filename); os.IsNotExist(err) {
		return nil // File does not exist, use default values.
	}

	data, err := ioutil.ReadFile(u.filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, u)
}
