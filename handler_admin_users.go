package main

import (
	"log"
	"net/http"
)

func (app *app) adminListUsers(w http.ResponseWriter, r *http.Request) {

	users, err := app.userService.GetAll()
	if err != nil {
		StatusCode(w, http.StatusInternalServerError)
		log.Panic(err)
	}
	serveTemplate(w, TemplateAdminListUsers, users)
}

func (app *app) adminAddUser(w http.ResponseWriter, r *http.Request) {

}

func (app *app) adminDeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (app *app) adminEditUser(w http.ResponseWriter, r *http.Request) {

}
