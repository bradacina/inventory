package admin

import (
	"log"
	"net/http"

	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/httphelp"
	"github.com/bradacina/inventory/routing"
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

type adminAddOrEditInventoryTemplateValues struct {
	*AdminEditInventoryRecord

	Action string
	Error  string
}

type AdminEditInventoryRecord struct {
	ID        int
	UserID    int
	Name      string
	Items     []AdminEditInventoryItemRecord
	IsDeleted bool
}

type AdminEditInventoryItemRecord struct {
	Quantity int
	SKU      string
	Title    string
	Barcode  string
}

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
}

func (ah *AdminHandler) ListInventories(w http.ResponseWriter, r *http.Request) {
	invs, err := ah.InventoryService.GetAll()
	if err != nil {
		log.Panicln(err)
		return
	}

	users, err := ah.UserService.GetAll()
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

func (ah *AdminHandler) AddInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		httphelp.ServeTemplate(w, httphelp.TemplateAdminInventory,
			adminAddOrEditInventoryTemplateValues{&AdminEditInventoryRecord{}, routing.RouteAdminAddInventory, ""})
		return
	} else if r.Method == http.MethodPost {
		inventory, err := parseInventoryFromRequest(r)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, routing.RouteAdminListInventories, http.StatusSeeOther)
			return
		}

		err = ah.InventoryService.CreateByAdmin(inventory)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, routing.RouteAdminListInventories, http.StatusSeeOther)
	}
}

func (ah *AdminHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	invID, err := httphelp.ParseIDFromQueryString(r)
	if err != nil {
		log.Println(err)
		httphelp.StatusCode(w, http.StatusInternalServerError)
		return
	}

	err = ah.InventoryService.SoftDeleteByAdmin(invID)
	if err != nil {
		log.Println("Error in adminDeleteInventory", err)
		httphelp.StatusCode(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, routing.RouteAdminListInventories, http.StatusSeeOther)
}

func (ah *AdminHandler) EditInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		invID, err := httphelp.ParseIDFromQueryString(r)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusInternalServerError)
			return
		}

		inventory, err := ah.InventoryService.GetByIDByAdmin(invID)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusNotFound)
		}

		invEdit, _ := mapToAdminEditInventoryRecord(inventory)

		httphelp.ServeTemplate(w, httphelp.TemplateAdminInventory,
			adminAddOrEditInventoryTemplateValues{invEdit, routing.RouteAdminEditInventory, ""})

	} else if r.Method == http.MethodPost {
		inv, err := parseInventoryFromRequest(r)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusInternalServerError)
			return
		}

		err = ah.InventoryService.UpdateByAdmin(inv)
		if err != nil {
			log.Println(err)
			httphelp.StatusCode(w, http.StatusNotFound)
			return
		}

		http.Redirect(w, r, routing.RouteAdminListInventories, http.StatusSeeOther)
	}
}

func parseInventoryFromRequest(r *http.Request) (*db.Inventory, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, err
	}

	var inv AdminEditInventoryRecord

	err = decoder.Decode(&inv, r.PostForm)

	if err != nil {
		return nil, err
	}

	inventory, err := mapToInventory(&inv)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func mapToInventory(inv *AdminEditInventoryRecord) (*db.Inventory, error) {
	inventory := db.Inventory{
		ID:        inv.ID,
		UserID:    inv.UserID,
		IsDeleted: inv.IsDeleted,
		Name:      inv.Name}

	for i := range inv.Items {
		inventory.Items = append(inventory.Items, db.Item(inv.Items[i]))
	}

	return &inventory, nil
}

func mapToAdminEditInventoryRecord(inv *db.Inventory) (*AdminEditInventoryRecord, error) {
	inventory := AdminEditInventoryRecord{
		ID:        inv.ID,
		UserID:    inv.UserID,
		IsDeleted: inv.IsDeleted,
		Name:      inv.Name}

	for _, el := range inv.Items {
		inventory.Items = append(inventory.Items, AdminEditInventoryItemRecord(el))
	}

	return &inventory, nil
}
