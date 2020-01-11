package main

import (
	"fmt"

	"github.com/orsenkucher/rave-era/repo"
)

func main() {
	r := repo.NewRepo()
	fmt.Println(r)
}
