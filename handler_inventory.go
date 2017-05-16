package main

import (
	"log"
	"net/http"
)

func (app *app) inventoryList(w http.ResponseWriter, r *http.Request) {
	user := GetLoggedInUser(r)
	if user == nil {
		log.Panic("User came back nil form request Context")
	}
	if r.Method == http.MethodGet {
		inventories := app.inventoryService.GetByUserID(user.ID)

		serveTemplate(w, TemplateInventoryList, inventories)

	} else if r.Method == http.MethodPost {
		inventories := app.inventoryService.GetByUserID(user.ID)

		//postInventories := parseInventories(r)

		// todo: determine deletions and make them as deleted

		// todo: determine additions

		// todo: determine name edits

		// todo: save in db

		// todo: serve new iventories through template

		serveTemplate(w, TemplateInventoryList, inventories)
	}
}
