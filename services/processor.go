package services

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Constants
const (
	Delay = 4.0 // Delay in minutes.
)

// TxProcessor processes the transaction list coming from the TronAPI service.
type TxProcessor struct {
	tronFetcher *TronFetcher
	carrier     *Carrier
	user        *User

	stopChans    sync.Map // Map of channels to signal stop processing for each wallet.
	isProcessing sync.Map // Map to store processing status for each wallet.
	lastTxIDs    sync.Map // Map to store the last transaction ID for each wallet.

	wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish.
}

// NewTxProcessor initializes a new transaction processor.
func NewTxProcessor(fetcher *TronFetcher, carrier *Carrier, user *User) *TxProcessor {
	return &TxProcessor{
		tronFetcher: fetcher,
		carrier:     carrier,
		user:        user,
	}
}

// StartProcessing starts processing transactions for all user wallets.
func (tx *TxProcessor) StartProcessing(bot *tgbotapi.BotAPI) {
	if len(tx.user.WalletAddresses) == 0 {
		message := "Notifications were not started due to a lack of wallets."
		log.Println("Processor: ", message)
		tx.carrier.SendToAdmin(bot, message)
		return
	}

	for _, wallet := range tx.user.WalletAddresses {
		// Check if already processing.
		if processing, ok := tx.isProcessing.Load(wallet); ok && processing.(bool) {
			message := fmt.Sprintf("The process is already underway for wallet ...%s. Start notifications call ignored.", wallet[len(wallet)-4:])
			log.Println("Processor: ", message)
			tx.carrier.SendToAdmin(bot, message)
			continue
		}

		// Set isProcessing to true.
		tx.isProcessing.Store(wallet, true)

		// Create a new stopChan if the existing one is closed.
		stopChanInterface, _ := tx.stopChans.LoadOrStore(wallet, make(chan struct{}))
		stopChan, ok := stopChanInterface.(chan struct{})
		if !ok || isClosed(stopChan) {
			stopChan = make(chan struct{})
			tx.stopChans.Store(wallet, stopChan)
		}

		message := fmt.Sprintf("Notifications were started for wallet ...%s.", wallet[len(wallet)-4:])
		log.Println("Processor: ", message)
		tx.carrier.SendToAdmin(bot, message)

		tx.wg.Add(1)
		go tx.processWallet(bot, wallet, stopChan)
	}
}

// processWallet processes transactions for a specific wallet.
func (tx *TxProcessor) processWallet(bot *tgbotapi.BotAPI, wallet string, stopChan chan struct{}) {
	defer tx.wg.Done() // Decrease the WaitGroup counter when the goroutine completes.

	for {
		select {
		case <-stopChan: // Check if a stop signal was sent.
			return
		default:
			// Fetch transactions for the given wallet.
			transactions, err := tx.tronFetcher.FetchTransactions(bot, wallet)
			if err != nil {
				log.Println("Processor: Error fetching transactions:", err)
				time.Sleep(time.Duration(Delay) * time.Minute) // Delay before the next request.
				continue
			}

			// Check the existence of transactions.
			if len(transactions) == 0 {
				log.Printf("Processor: No transactions found for wallet ...%s", wallet[len(wallet)-4:])
				time.Sleep(time.Duration(Delay) * time.Minute) // Delay before the next request.
				continue
			}

			// Get the latest transaction.
			lastTransaction := transactions[0]
			tx.lastTxIDs.LoadOrStore(wallet, "")

			// If the latest transaction is the same as the previous one, wait and check later.
			if lastTransaction.TransactionID == tx.getLastTransactionID(wallet) {
				time.Sleep(time.Duration(Delay) * time.Minute) // Delay before the next request.
				continue
			}

			// Update the latest transaction.
			tx.lastTxIDs.Store(wallet, lastTransaction.TransactionID)

			// Check if the transaction is old.
			if tx.isTransactionOld(lastTransaction.Timestamp) {
				time.Sleep(time.Duration(Delay) * time.Minute) // Delay before the next request.
				continue
			}

			// Remove trailing zeros from the transaction value.
			valueInUSDT := tx.convertValueToUSDT(lastTransaction.Value)
			message := fmt.Sprintf("At %s, a transfer of %.5f USDT was made.\nRecipient wallet address: ...%s",
				time.Unix(lastTransaction.Timestamp/1000, 0).Format("2006-01-02 15:04:05"), valueInUSDT, wallet[len(wallet)-4:])
			log.Println("Processor:", message)

			// Send the message to all recipients and the administrator.
			tx.carrier.SendToAll(bot, tx.user.Recipients, message)
			tx.carrier.SendToAdmin(bot, message)

			time.Sleep(time.Duration(Delay) * time.Minute) // Delay before the next request.
		}
	}
}

// StopProcessing stops processing transactions for all user wallets.
func (tx *TxProcessor) StopProcessing(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	tx.stopChans.Range(func(key, value interface{}) bool {
		if stopChan, ok := value.(chan struct{}); ok && !isClosed(stopChan) {
			close(stopChan)
		}

		tx.isProcessing.Store(key, false)
		tx.stopChans.Delete(key) // Remove the stopChan after stopping.

		return true
	})

	tx.wg.Wait()

	message := "Notifications were stopped."
	log.Println("Processor: ", message)
	tx.carrier.SendToAdmin(bot, message)
}

// CheckStatus checks the processing status of all user wallets and sends it to the admin.
func (tx *TxProcessor) CheckStatus(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	statusMessage := fmt.Sprintf("<b>Recipients:</b> %v\n\n", tx.user.Recipients)

	for i, wallet := range tx.user.WalletAddresses {
		processing, exists := tx.isProcessing.Load(wallet)
		if !exists {
			processing = false
		}

		lastTxID, exists := tx.lastTxIDs.Load(wallet)
		if !exists {
			lastTxID = "N/A"
		}

		statusMessage += fmt.Sprintf("<b>Wallet %d:</b> ...%s\nisProcessing: %t\nLast transaction ID: %s\n\n", i+1, wallet[len(wallet)-4:], processing, lastTxID)
	}

	tx.carrier.SendToAdmin(bot, statusMessage)
}

// Helper functions

// getLastTransactionID returns the last transaction ID for a given wallet.
func (tx *TxProcessor) getLastTransactionID(wallet string) string {
	if lastTxID, ok := tx.lastTxIDs.Load(wallet); ok {
		return lastTxID.(string)
	}

	return ""
}

// isTransactionOld checks if a transaction is old based on its timestamp.
func (tx *TxProcessor) isTransactionOld(timestamp int64) bool {
	transactionTime := time.Unix(timestamp/1000, 0)
	return time.Since(transactionTime) > time.Duration(1.5*Delay)*time.Minute
}

// convertValueToUSDT removes trailing zeros from the transaction value and formats it to 2 decimal places.
func (tx *TxProcessor) convertValueToUSDT(value string) float64 {
	v, _ := strconv.Atoi(value)
	return math.Round(float64(v)/10000) / 100
}
