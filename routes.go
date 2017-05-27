package main

import (
	"log"
	"net/http"

	"github.com/bradacina/inventory/httphelp"
)

func registerRoutes(app *app) {
	log.Println("Setting up routes...")
	http.HandleFunc("/register", app.register)
	http.HandleFunc("/login", app.login)
	http.HandleFunc("/logout", app.logout)

	http.HandleFunc("/admin_users",
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetHandleFunc(
					app.adminListUsers))))

	http.HandleFunc("/admin_delete_user",
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.PostHandleFunc(
					app.adminDeleteUser))))

	http.HandleFunc("/admin_edit_user",
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.PostHandleFunc(
					app.adminEditUser))))

	http.HandleFunc("/admin_list_inventories",
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetHandleFunc(
					app.adminListInventories))))

	http.HandleFunc("/admin_add_inventory",
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetPostHandleFunc(
					app.adminAddInventory))))

	http.HandleFunc("/secure", httphelp.AuthenticatedHandleFunc(app, app.secure))
	log.Println("Done")
}
