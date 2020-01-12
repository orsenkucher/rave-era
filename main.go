package main

import (
	"fmt"

	"github.com/orsenkucher/rave-era/repo"
)

func main() {
	r := repo.NewRepo()
	var k int64
	k = 0
	raver := repo.Raver{Name: "Srgey",
		LastName: "Cheremshinsky",
		ID:       k,
		Uni:      "KNU",
		Age:      18,
	}
	fmt.Println(r.Events)
	r.AddUser(raver)
	r.Subscribe(raver, []repo.Raver{raver}, "Новогодний рейв")
	//b := bot.NewBot(creds.CrRvra, r)
	//fmt.Println("Listening..")
	//b.Listen()
}
