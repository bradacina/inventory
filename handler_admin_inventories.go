package main

import (
	"log"
	"net/http"

	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/httphelp"
	"github.com/gorilla/schema"
)

type adminInventoryListRecord struct {
	ID         int
	Name       string
	UserID     int
	Email      string
	TotalItems int
	IsDeleted  bool
}

type adminViewInventoryEdit struct {
	*AdminInventoryEdit
	Action string
}

type AdminInventoryEdit struct {
	ID        int
	UserID    int
	Name      string
	Items     []AdminInventoryItemEdit
	IsDeleted bool
}

type AdminInventoryItemEdit struct {
	Quantity int
	SKU      string
	Title    string
	Barcode  string
}

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
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

		// todo: improve this with a map[][]
		for j := range users {
			if users[j].ID == invs[i].UserID {
				record.Email = users[j].Email
				break
			}
		}

		displayInvs = append(displayInvs, record)
	}
	httphelp.ServeTemplate(w, httphelp.TemplateAdminListInventories, displayInvs)
}

func (app *app) adminAddInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		httphelp.ServeTemplate(w, httphelp.TemplateAdminInventory,
			adminViewInventoryEdit{&AdminInventoryEdit{}, "/admin_add_inventory"})
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

func (app *app) adminDeleteInventory(w http.ResponseWriter, r *http.Request) {
	invID, err := httphelp.ParseIDFromQueryString(r)
	if err != nil {
		log.Println(err)
		httphelp.StatusCode(w, http.StatusInternalServerError)
		return
	}

	err = app.inventoryService.SoftDeleteByAdmin(invID)
	if err != nil {
		log.Println("Error in adminDeleteInventory", err)
		httphelp.StatusCode(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin_list_inventories", http.StatusSeeOther)
}

func (app *app) adminEditInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		invID, err := httphelp.ParseIDFromQueryString(r)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusInternalServerError)
			return
		}

		inventory, err := app.inventoryService.GetByIDByAdmin(invID)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusNotFound)
		}

		invEdit, _ := mapToAdminInventoryEdit(inventory)

		httphelp.ServeTemplate(w, httphelp.TemplateAdminInventory,
			adminViewInventoryEdit{invEdit, "/admin_edit_inventory"})

	} else if r.Method == http.MethodPost {

	}
}

func parseInventoryFromRequest(r *http.Request) (*db.Inventory, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, err
	}

	var inv AdminInventoryEdit

	err = decoder.Decode(&inv, r.PostForm)

	if err != nil {
		return nil, err
	}

	log.Printf("%v\r\n", inv)

	inventory, err := mapToInventory(&inv)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func mapToInventory(inv *AdminInventoryEdit) (*db.Inventory, error) {
	inventory := db.Inventory{
		ID:        inv.ID,
		UserID:    inv.UserID,
		IsDeleted: inv.IsDeleted,
		Name:      inv.Name}

	for i := range inv.Items {
		inventory.Items = append(inventory.Items, db.Item{
			Barcode:  inv.Items[i].Barcode,
			Quantity: inv.Items[i].Quantity,
			SKU:      inv.Items[i].SKU,
			Title:    inv.Items[i].Title})
	}

	return &inventory, nil
}

func mapToAdminInventoryEdit(inv *db.Inventory) (*AdminInventoryEdit, error) {
	inventory := AdminInventoryEdit{
		ID:        inv.ID,
		UserID:    inv.UserID,
		IsDeleted: inv.IsDeleted,
		Name:      inv.Name}

	for _, el := range inv.Items {
		inventory.Items = append(inventory.Items, AdminInventoryItemEdit(el))
	}

	return &inventory, nil
}
