package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/rave-era/creds"
	"github.com/orsenkucher/rave-era/repo"
)

type Bot struct {
	cred creds.Credential
	repo *repo.Repo
	api  *tgbotapi.BotAPI
}

func NewBot(cred creds.Credential, repo *repo.Repo) *Bot {
	b := &Bot{cred: cred, repo: repo}
	b.initAPI()
	return b
}

func (b *Bot) initAPI() {
	var err error
	b.api, err = tgbotapi.NewBotAPI(b.cred.String())
	if err != nil {
		log.Panic(err)
	}

	b.api.Debug = false
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	_, err = b.api.RemoveWebhook()
	if err != nil {
		log.Println("Cant remove webhook")
	}
}

func (b *Bot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleCallback(update)
			continue
		}

		if update.Message != nil {
			if update.Message.IsCommand() {
				b.handleCommand(update)
			} else if update.Message.Text != "" {
				b.handleMessage(update)
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			continue
		}
	}
}
