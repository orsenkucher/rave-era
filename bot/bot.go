type Bot struct {
	repo *Repo
}

func NewBot(repo *Repo) *Bot {
	return &Bot{repo: repo}
}

func (b *Bot) Listen() {

}
