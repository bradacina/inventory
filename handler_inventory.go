package main

import (
	"log"
	"net/http"
	"strconv"
)

func (app *app) inventoryList(w http.ResponseWriter, r *http.Request) {
	user := GetLoggedInUser(r)
	if user == nil {
		log.Panic("User came back nil form request Context")
	}

	inventories, err := app.inventoryService.GetByUserID(user.ID)
	if err != nil {
		log.Panic(err)
	}

	if r.Method == http.MethodGet {
		serveTemplate(w, TemplateInventoryList, inventories)

	} else if r.Method == http.MethodPost {
		err := app.inventoryService.CreateWithName(r.FormValue("name"), user.ID)
		if err != nil {
			log.Println(err)
			serveTemplate(w, TemplateInventoryList, inventories)
			return
		}

		inventories, err := app.inventoryService.GetByUserID(user.ID)
		if err != nil {
			log.Panic(err)
		}

		serveTemplate(w, TemplateInventoryList, inventories)
	} else if r.Method == http.MethodDelete {

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Println(err)
			serveTemplate(w, TemplateInventoryList, inventories)
			return
		}

		err = app.inventoryService.SoftDelete(id, user.ID)
		if err != nil {
			log.Println(err)
		}

		serveTemplate(w, TemplateInventoryList, inventories)
	}
}
