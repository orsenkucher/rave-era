package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/rave-era/repo"
)

// Item represens context of work
type Item struct {
	DestID  int64  // Item destination
	CallID  string // Callback
	Payload string
	OnMsg   func(text string)
}

func (b *Bot) handleCallback(update tgbotapi.Update) {
	fmt.Println("Callback:", update)
	data := update.CallbackQuery.Data
	destID := update.CallbackQuery.Message.Chat.ID
	callID := update.CallbackQuery.ID
	payload := strings.Split(data, ":")[1]
	switch {
	case strings.Contains(data, "free"):
		b.onFree(Item{
			DestID:  destID,
			CallID:  callID,
			Payload: payload,
		})
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	fmt.Println("Message:", update)
	msg := update.Message
	destID := msg.Chat.ID
	text := msg.Text
	if it, ok := b.jobs[destID]; ok && it.OnMsg != nil {
		it.OnMsg(text)
		it.OnMsg = nil
	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	cmd := update.Message.Command()
	fmt.Println("Command:", cmd)
	switch cmd {
	case "start", "go", "rave":
		it := Item{DestID: update.Message.Chat.ID}
		b.jobs[it.DestID] = &it
		b.onRave(it)
	default:
		return
	}
}

func (b *Bot) onFree(it Item) {
	fmt.Println("FREE")
	// raver := repo.Raver{ID: it.DestID}
	raver, ok := b.repo.Find(it.DestID)
	if !ok {
		log.Println("Error: User not found in repo")
	}
	friends := []repo.Raver{}
	b.repo.Subscribe(raver, friends, it.Payload)
	snackMsg := fmt.Sprintf("🔥%s: %s+%v", it.Payload, raver.Name, len(friends))
	_, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(it.CallID, snackMsg))
	if err != nil {
		log.Println(err)
	}
	msg := tgbotapi.NewMessage(it.DestID, snackMsg)
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) onRave(it Item) {
	fmt.Println("RAVIN'")
	if !b.repo.UserRegistrated(it.DestID) {
		fmt.Printf("Raver with ID %v was not regestered\n", it.DestID)
		b.onReg(it)
		return
	}
	events := b.repo.Events
	fmt.Println(events)
	text := "🔥RAVE LIST"
	msg := tgbotapi.NewMessage(it.DestID, text)
	if len(events) > 0 {
		mkp, _ := genFreeList(events)
		msg.ReplyMarkup = &mkp
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) onReg(it Item) {
	b.jobs[it.DestID].OnMsg = func(text string) {
		triple := strings.Split(text, " ")
		if len(triple) != 3 { // TODO
			msg := tgbotapi.NewMessage(it.DestID, "try again")
			if _, err := b.api.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}
		name := triple[0]
		lastname := triple[1]
		age, err := strconv.Atoi(triple[2])
		if err != nil { // TODO
			msg := tgbotapi.NewMessage(it.DestID, "try again")
			if _, err := b.api.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}
		raver := repo.Raver{
			ID:       it.DestID,
			Name:     name,
			LastName: lastname,
			Age:      age,
		}
		b.repo.AddUser(raver)
		b.onRave(it)
	}
	text := `🎊 RAVER REGISTRATION 🎊
Enter your name and age
ex: 'Orsen Kucher 20'`
	msg := tgbotapi.NewMessage(it.DestID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}
