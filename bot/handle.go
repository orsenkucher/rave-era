package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleCallback(update tgbotapi.Update) {
	fmt.Println("Callback:", update)
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	fmt.Println("Command:", update)
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	fmt.Println("Message:", update)
}
