package main

import (
	"flag"
	"log"

	"os"

	"github.com/asdine/storm"
	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/services"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {

	dbpath := flag.String("d", "", "Path to the `database`")
	userId := flag.Int("u", 0, "`UserId` to underlete")

	flag.Parse()

	_, err := os.Stat(*dbpath)
	if err != nil {
		log.Fatal(err)
	}

	if len(*dbpath) == 0 || *userId == 0 || flag.NFlag() != 2 {
		flag.PrintDefaults()
	}

	database, err := storm.Open(*dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	userRepo := db.NewUserRepo(database)

	userService := services.NewUserService(userRepo)

	//inventoryRepo := db.NewInventoryRepo(database)
	//inventoryService := services.NewInventoryService(inventoryRepo, userRepo)

	user, err := userService.GetByID(*userId)
	if err != nil {
		log.Fatal(err)
	}

	user.IsDeleted = false

	err = userService.UpdateByAdmin(user)
	if err != nil {
		log.Fatal(err)
	}

}
