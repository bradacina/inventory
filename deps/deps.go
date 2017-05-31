package deps

import (
	"log"

	"github.com/asdine/storm"
	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/logincookie"
	"github.com/bradacina/inventory/routing"
	"github.com/bradacina/inventory/services"
)

type Deps struct {
	Db               *storm.DB
	UserRepo         db.UserRepoer
	InventoryRepo    db.InventoryRepoer
	UserService      services.UserServicer
	InventoryService services.InventoryServicer
	CookieAuth       *logincookie.CookieAuthentication
}

func NewDeps() *Deps {
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

	return &Deps{
		database,
		userRepo,
		inventoryRepo,
		userService,
		inventoryService,
		cookieAuth}
}

func (deps *Deps) GetCookieAuthentication() *logincookie.CookieAuthentication {
	return deps.CookieAuth
}

func (deps *Deps) GetUserService() services.UserServicer {
	return deps.UserService
}

func (deps *Deps) GetLoginRoute() string {
	return routing.RouteLogin
}
