package main

import (
	"log"
	"net/http"
)

func registerRoutes(app *app) {
	log.Println("Setting up routes...")
	http.HandleFunc("/register", app.register)
	http.HandleFunc("/login", app.login)
	http.HandleFunc("/logout", app.logout)

	http.HandleFunc("/admin_users",
		app.authenticatedHandleFunc(
			app.isAdminHandleFunc(
				app.getHandleFunc(
					app.adminListUsers))))

	http.HandleFunc("/admin_add_user",
		app.authenticatedHandleFunc(
			app.isAdminHandleFunc(
				app.postHandleFunc(
					app.adminAddUser))))

	http.HandleFunc("/admin_delete_user",
		app.authenticatedHandleFunc(
			app.isAdminHandleFunc(
				app.postHandleFunc(
					app.adminDeleteUser))))

	http.HandleFunc("/admin_edit_user",
		app.authenticatedHandleFunc(
			app.isAdminHandleFunc(
				app.postHandleFunc(
					app.adminEditUser))))

	http.HandleFunc("/secure", app.authenticatedHandleFunc(app.secure))
	log.Println("Done")
}
