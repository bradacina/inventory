package main

import (
	"log"
	"net/http"

	"github.com/bradacina/inventory/httphelp"
)

const (
	RouteRegister             = "/register"
	RouteLogin                = "/login"
	RouteLogout               = "/logout"
	RouteAdminListUsers       = "/admin_list_users"
	RouteAdminDeleteUser      = "/admin_delete_user"
	RouteAdminEditUser        = "/admin_edit_user"
	RouteAdminListInventories = "/admin_list_inventories"
	RouteAdminEditInventory   = "/admin_edit_inventory"
	RouteAdminDeleteInventory = "/admin_delete_inventory"
	RouteAdminAddInventory    = "/admin_add_inventory"
)

func registerRoutes(app *app) {
	log.Println("Setting up routes...")
	http.HandleFunc(RouteRegister, app.register)
	http.HandleFunc(RouteLogin, app.login)
	http.HandleFunc(RouteLogout, app.logout)

	http.HandleFunc(RouteAdminListUsers,
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetHandleFunc(
					app.adminListUsers))))

	http.HandleFunc(RouteAdminDeleteUser,
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.PostHandleFunc(
					app.adminDeleteUser))))

	http.HandleFunc(RouteAdminEditUser,
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetPostHandleFunc(
					app.adminEditUser))))

	http.HandleFunc(RouteAdminListInventories,
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetHandleFunc(
					app.adminListInventories))))

	http.HandleFunc(RouteAdminAddInventory,
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetPostHandleFunc(
					app.adminAddInventory))))

	http.HandleFunc(RouteAdminDeleteInventory,
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.PostHandleFunc(
					app.adminDeleteInventory))))

	http.HandleFunc(RouteAdminEditInventory,
		httphelp.AuthenticatedHandleFunc(
			app,
			httphelp.IsAdminHandleFunc(
				httphelp.GetPostHandleFunc(
					app.adminEditInventory))))

	http.HandleFunc("/secure", httphelp.AuthenticatedHandleFunc(app, app.secure))
	log.Println("Done")
}
