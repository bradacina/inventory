package main

import (
	"log"
	"net/http"

	"github.com/asdine/storm"
	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/logincookie"
	"github.com/bradacina/inventory/services"
)

type app struct {
	db               *storm.DB
	userRepo         db.UserRepoer
	inventoryRepo    db.InventoryRepoer
	userService      services.UserServicer
	inventoryService services.InventoryServicer
	cookieAuth       *logincookie.CookieAuthentication
}

func newApp() *app {
	database, err := storm.Open("inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := db.NewUserRepo(database)
	inventoryRepo := db.NewInventoryRepo(database)

	userService := services.NewUserService(userRepo)

	inventoryService := services.NewInventoryService(inventoryRepo, userRepo)

	cookieAuth := logincookie.NewCookieAuthentication("test.dev",
		[]byte("what a secret123"),
		[]byte("silent night all"))

	return &app{
		database,
		userRepo,
		inventoryRepo,
		userService,
		inventoryService,
		cookieAuth}
}

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	app := newApp()
	registerRoutes(app)

	log.Println("Listening at https://localhost:9080")
	log.Fatal(http.ListenAndServeTLS(":9080", "cert/pub.key", "cert/priv.key", nil))

	app.db.Close()
}
