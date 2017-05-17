package main

import (
	"log"
	"net/http"

	"github.com/asdine/storm"
)

type app struct {
	db               *storm.DB
	userRepo         UserRepoer
	inventoryRepo    InventoryRepoer
	userService      UserServicer
	inventoryService InventoryServicer
	cookieHelper     *cookieHelper
}

func newApp() *app {
	db, err := storm.Open("inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := newUserRepo(db)
	inventoryRepo := newInventoryRepo(db)

	userService := NewUserService(userRepo)

	inventoryService := NewInventoryService(inventoryRepo)

	cookieHelper := cookieHelper{"test.dev",
		[]byte("what a secret123"),
		[]byte("silent night all")}

	return &app{
		db,
		userRepo,
		inventoryRepo,
		userService,
		inventoryService,
		&cookieHelper}
}

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	app := newApp()
	registerRoutes(app)

	log.Println("Listening at https://localhost:9080")
	log.Fatal(http.ListenAndServeTLS(":9080", "cert/pub.key", "cert/priv.key", nil))

	app.db.Close()
}
