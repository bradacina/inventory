package main

import (
	"log"
	"net/http"

	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/httphelp"
)

type adminInventoryListRecord struct {
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

	var displayInvs []adminInventoryListRecord
	for i := range invs {
		record := adminInventoryListRecord{
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

	httphelp.ServeTemplate(w, httphelp.TemplateAdminListInventories, displayInvs)
}

func (app *app) adminAddInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		httphelp.ServeTemplate(w, httphelp.TemplateAdminInventory, db.Inventory{})
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

func parseInventoryFromRequest(r *http.Request) (*db.Inventory, error) {

}
