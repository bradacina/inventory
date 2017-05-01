package main

import (
	"log"

	"github.com/asdine/storm"
)

func main() {
	dbMain, err := storm.Open("c:\\temp\\storm.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbMain.Close()

	userRepo := newUserRepo(dbMain)
	invRepo := newInventoryRepo(dbMain)

	user, _ := userRepo.GetByEmail("bogdan@email.com")
	invs := invRepo.GetAll()

	log.Print(user, invs)
}
