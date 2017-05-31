package inventory

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bradacina/inventory/deps"
	"github.com/bradacina/inventory/httphelp"
)

type InventoryHandler struct {
	*deps.Deps
}

func NewHandler(deps *deps.Deps) *InventoryHandler {
	return &InventoryHandler{deps}
}

func (ih *InventoryHandler) InventoryList(w http.ResponseWriter, r *http.Request) {
	user := httphelp.GetLoggedInUser(r)
	if user == nil {
		log.Panic("User came back nil form request Context")
	}

	inventories, err := ih.InventoryService.GetByUserID(user.ID)
	if err != nil {
		log.Panic(err)
	}

	if r.Method == http.MethodGet {
		httphelp.ServeTemplate(w, httphelp.TemplateInventoryList, inventories)

	} else if r.Method == http.MethodPost {
		err := ih.InventoryService.CreateWithName(r.FormValue("name"), user.ID)
		if err != nil {
			log.Println(err)
			httphelp.ServeTemplate(w, httphelp.TemplateInventoryList, inventories)
			return
		}

		inventories, err := ih.InventoryService.GetByUserID(user.ID)
		if err != nil {
			log.Panic(err)
		}

		httphelp.ServeTemplate(w, httphelp.TemplateInventoryList, inventories)
	} else if r.Method == http.MethodDelete {

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Println(err)
			httphelp.ServeTemplate(w, httphelp.TemplateInventoryList, inventories)
			return
		}

		err = ih.InventoryService.SoftDelete(id, user.ID)
		if err != nil {
			log.Println(err)
		}

		httphelp.ServeTemplate(w, httphelp.TemplateInventoryList, inventories)
	}
}
