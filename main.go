package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

const botApp = "6766212541:AAGQTAPSbZGQoJ-FmG-nyYZLPCQQCc5wFIw"

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
		break
	}
}

var users = make(map[int64]int64)

func handleMessage(message *tgbotapi.Message) {
	// Print to console
	log.Printf("message %+v", message)

	users[message.From.ID] = message.Chat.ID

	fmt.Println(message.From.ID, "   equal   ", message.Chat.ID)
}

func init() {
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
