package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Item represens context of work
type Item struct {
	DestID int64 // Item destination
}

func (b *Bot) handleCallback(update tgbotapi.Update) {
	fmt.Println("Callback:", update)
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	cmd := update.Message.Command()
	fmt.Println("Command:", cmd)
	switch cmd {
	case "start", "go", "rave":
		b.onRave(Item{
			DestID: update.Message.Chat.ID,
		})
	default:
		return
	}
}

func (b *Bot) onRave(it Item) {
	fmt.Println("RAVIN'")
	events := b.repo.LoadEvents()
	fmt.Println(events)
	msg := tgbotapi.NewMessage(it.DestID, "RAVE LIST")
	if len(events) > 0 {
		mkp, _ := genFreeList(events)
		msg.ReplyMarkup = &mkp
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	fmt.Println("Message:", update)
}
