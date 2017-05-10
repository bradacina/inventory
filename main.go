package main

import (
	"log"
	"net/http"

	"github.com/asdine/storm"
)

type app struct {
	db            *storm.DB
	userRepo      UserRepoer
	inventoryRepo InventoryRepoer
	userService   UserServicer
}

func newApp() *app {
	db, err := storm.Open("test.db")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := newUserRepo(db)
	inventoryRepo := newInventoryRepo(db)

	userService := NewUserService(userRepo)

	return &app{
		db,
		userRepo,
		inventoryRepo,
		userService}
}

func main() {
	app := newApp()
	registerRoutes(app)

	log.Println("Listening at http://localhost:9080")
	log.Fatal(http.ListenAndServe(":9080", nil))

	app.db.Close()
}
