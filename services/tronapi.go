package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Service for receiving data on cryptotransactions.
type TronFetcher struct {
	carrier         *Carrier
	tronApiEndpoint string
}

func NewTronFetcher(carrier *Carrier, tronApiEndpoint string) *TronFetcher {
	return &TronFetcher{
		carrier:         carrier,
		tronApiEndpoint: tronApiEndpoint,
	}
}

type Transaction struct {
	Timestamp     int64  `json:"block_timestamp"`
	Value         string `json:"value"`
	TransactionID string `json:"transaction_id"`
}

type Response struct {
	Data []Transaction `json:"data"`
}

// Get transactions of a specific cryptocurrency wallet.
func (t *TronFetcher) FetchTransactions(bot *tgbotapi.BotAPI, walletAddress string) ([]Transaction, error) {
	// Check if the tracking address is defined.
	if walletAddress == "" {
		message := "Address is not defined."
		log.Println("TronAPI: ", message)
		t.carrier.SendToAdmin(bot, message)

		return nil, errors.New("address is not defined")
	}

	// Connect to the API.
	url := fmt.Sprintf(t.tronApiEndpoint, walletAddress)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("TronAPI: Error trying to fetch transactions: %v\n", err)
		return nil, err
	}
	defer response.Body.Close()

	// Parse the incoming data.
	var respData Response
	if err := json.NewDecoder(response.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return respData.Data, nil
}
