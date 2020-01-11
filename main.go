package main

import (
	"fmt"

	"github.com/orsenkucher/rave-era/bot"
	"github.com/orsenkucher/rave-era/creds"
	"github.com/orsenkucher/rave-era/repo"
)

func main() {
	r := repo.NewRepo()
	b := bot.NewBot(creds.CrRvra, r)
	fmt.Println("Listening..")
	b.Listen()
}
