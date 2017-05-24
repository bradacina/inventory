package main

import (
	"log"
	"net/http"
	"strconv"
)

func (app *app) adminListUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.userService.GetAll()
	if err != nil {
		StatusCode(w, http.StatusInternalServerError)
		log.Panic(err)
	}
	serveTemplate(w, TemplateAdminListUsers, users)
}

func (app *app) adminDeleteUser(w http.ResponseWriter, r *http.Request) {
	requestUserID := r.FormValue("userID")
	userID, err := strconv.Atoi(requestUserID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin_list_users", http.StatusSeeOther)
		return
	}

	user, err := app.userService.GetByID(userID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/admin_list_users", http.StatusSeeOther)
		return
	}

	user.IsDeleted = true
	err = app.userService.UpdateByAdmin(user)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/admin_list_users", http.StatusSeeOther)
}

func (app *app) adminEditUser(w http.ResponseWriter, r *http.Request) {
	// todo: implement
	http.Redirect(w, r, "/admin_list_users", http.StatusSeeOther)
}
