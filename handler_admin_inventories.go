package main

import (
	"log"
	"net/http"
)

type inventoryListRecord struct {
	ID         int
	Name       string
	UserID     int
	Email      string
	TotalItems int
	IsDeleted  bool
}

func (app *app) adminListInventories(w http.ResponseWriter, r *http.Request) {
	invs, err := app.inventoryService.GetAll()
	if err != nil {
		log.Panicln(err)
		return
	}

	users, err := app.userService.GetAll()
	if err != nil {
		log.Println(err)
	}

	var displayInvs []inventoryListRecord
	for i := range invs {
		record := inventoryListRecord{
			ID:         invs[i].ID,
			UserID:     invs[i].UserID,
			Name:       invs[i].Name,
			IsDeleted:  invs[i].IsDeleted,
			TotalItems: len(invs[i].Items)}

		for j := range users {
			if users[j].ID == invs[i].UserID {
				record.Email = users[i].Email
				break
			}
		}

		displayInvs = append(displayInvs, record)
	}

	serveTemplate(w, TemplateAdminListInventories, displayInvs)
}

func (app *app) adminAddInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		serveTemplate(w, TemplateAdminInventory, Inventory{})
		return
	} else if r.Method == http.MethodPost {
		inventory, err := parseInventoryFromRequest(r)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/admin_list_inventories", http.StatusSeeOther)
			return
		}

		err = app.inventoryService.CreateByAdmin(inventory)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/admin_list_inventories", http.StatusSeeOther)
	}
}

func parseInventoryFromRequest(r *http.Request) (*Inventory, error) {

}
