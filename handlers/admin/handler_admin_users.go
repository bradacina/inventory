package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bradacina/inventory/httphelp"
	"github.com/bradacina/inventory/routing"
)

func (ah *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {

	users, err := ah.UserService.GetAll()
	if err != nil {
		httphelp.StatusCode(w, http.StatusInternalServerError)
		log.Panic(err)
	}
	httphelp.ServeTemplate(w, httphelp.TemplateAdminListUsers, users)
}

func (ah *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	requestUserID := r.FormValue("userID")
	userID, err := strconv.Atoi(requestUserID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
		return
	}

	user, err := ah.UserService.GetByID(userID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
		return
	}

	user.IsDeleted = true
	err = ah.UserService.UpdateByAdmin(user)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
}

func (ah *AdminHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	// todo: implement
	http.Redirect(w, r, routing.RouteAdminListUsers, http.StatusSeeOther)
}
