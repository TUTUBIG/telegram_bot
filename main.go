package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"time"
)

const botApp = "6766212541:AAGQTAPSbZGQoJ-FmG-nyYZLPCQQCc5wFIw" // wali

var bot *tgbotapi.BotAPI

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI(botApp)
	if err != nil {
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// `updates` is a golang channel that receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	ctx, cancel := context.WithCancel(context.Background())
	receiveUpdates(ctx, updates)
	// broadcast()
	cancel()
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is canceled
		case <-ctx.Done():
			return
		// receive an update from a channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
	case update.PreCheckoutQuery != nil:
		handlePreCheckoutQuery(update.PreCheckoutQuery)
	case update.InlineQuery != nil:
		handleInlineQuery(update.InlineQuery)
	}
}

var users = make(map[int64]int64)

func handlePreCheckoutQuery(query *tgbotapi.PreCheckoutQuery) {
	log.Printf("query %+v", query)
	answer := tgbotapi.PreCheckoutConfig{
		PreCheckoutQueryID: query.ID,
		OK:                 true,
		ErrorMessage:       "",
	}
	if rand.Int31n(3) == 0 {
		answer.OK = false
		answer.ErrorMessage = "Please be patient, what you purchased for is in preparing!"
	}
	response, err := bot.Request(answer)
	if err != nil {
		log.Println("error: ", err.Error())
		return
	}
	if !response.Ok {
		log.Printf("error: %+v\n", response)
		return
	}
}

func handleMessage(message *tgbotapi.Message) {
	// Print to console
	log.Printf("message %+v", message)

	users[message.From.ID] = message.Chat.ID

	fmt.Println(message.From.ID, "   equal   ", message.Chat.ID)

	if message.IsCommand() {
		switch message.Command() {
		case "shop":
			demoInvoice := tgbotapi.NewInvoice(message.Chat.ID, "Digital Cars Set", "This is a demonstration for telegram payment with fait currency", "sku=100", "2051251535:TEST:OTk5MDA4ODgxLTAwNQ", "", "USD", []tgbotapi.LabeledPrice{
				{
					Label:  "1 * 🏎",
					Amount: 111,
				},
				{
					Label:  "5 * 🚗",
					Amount: 352,
				},
			})
			demoInvoice.PhotoURL = "https://pub-6c52100fa9ac41f681f0713eac878541.r2.dev/Aave.png"
			demoInvoice.MaxTipAmount = 50
			demoInvoice.SuggestedTipAmounts = []int{5}

			_, err := bot.Send(demoInvoice)
			if err != nil {
				log.Println("error: ", err.Error())
				return
			}

			demoInvoice1 := tgbotapi.NewInvoice(message.Chat.ID, "Digital Yacht", "This is a demonstration for telegram payment with telegram star", "sku=101", "", "", "XTR", []tgbotapi.LabeledPrice{
				{
					Label:  "1 * 🛥️",
					Amount: 1000,
				},
			})
			demoInvoice1.PhotoURL = "https://pub-6c52100fa9ac41f681f0713eac878541.r2.dev/Aave.png"
			demoInvoice1.SuggestedTipAmounts = []int{}
			bot.Debug = true
			_, err = bot.Send(demoInvoice1)
			if err != nil {
				log.Println("error: ", err.Error())
				return
			}
		}
	}

	if message.SuccessfulPayment != nil {
		payment := message.SuccessfulPayment
		log.Printf("payment id %s  %+v", payment.TelegramPaymentChargeID, payment)
	}
}

func handleInlineQuery(query *tgbotapi.InlineQuery) {
	// Print to console
	log.Printf("query %+v", query)

}

func broadcast() {
	go func() {
		for {
			time.Sleep(5 * time.Second)

			message := tgbotapi.NewMessage(1632669575, "hi")
			reply, err := bot.Send(message)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(reply)

		}
	}()
}
