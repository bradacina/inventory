package main

import (
	"log"
	"net/http"

	"github.com/bradacina/inventory/deps"
	"github.com/bradacina/inventory/handlers/admin"
	"github.com/bradacina/inventory/handlers/login"
	"github.com/bradacina/inventory/handlers/register"
	"github.com/bradacina/inventory/httphelp"
	"github.com/bradacina/inventory/routing"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	deps := deps.NewDeps()
	registerRoutes(deps)

	log.Println("Listening at https://127.0.0.1:9080")
	log.Fatal(http.ListenAndServeTLS(":9080", "cert/pub.key", "cert/priv.key", nil))

	deps.Db.Close()
}

func registerRoutes(deps *deps.Deps) {

	registerHandler := register.NewHandler(deps)
	loginHandler := login.NewHandler(deps)
	adminHandler := admin.NewHandler(deps)

	log.Println("Setting up routes...")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc(routing.RouteRegister, registerHandler.Register)
	http.HandleFunc(routing.RouteLogin, loginHandler.Login)
	http.HandleFunc(routing.RouteLogout, loginHandler.Logout)

	http.HandleFunc(routing.RouteAdminListUsers,
		httphelp.AuthenticatedHandleFunc(
			deps,
			httphelp.IsAdminHandleFunc(
				httphelp.GetHandleFunc(
					adminHandler.ListUsers))))

	http.HandleFunc(routing.RouteAdminDeleteUser,
		httphelp.AuthenticatedHandleFunc(
			deps,
			httphelp.IsAdminHandleFunc(
				httphelp.PostHandleFunc(
					adminHandler.DeleteUser))))

	http.HandleFunc(routing.RouteAdminEditUser,
		httphelp.AuthenticatedHandleFunc(
			deps,
			httphelp.IsAdminHandleFunc(
				httphelp.GetPostHandleFunc(
					adminHandler.EditUser))))

	http.HandleFunc(routing.RouteAdminListInventories,
		httphelp.AuthenticatedHandleFunc(
			deps,
			httphelp.IsAdminHandleFunc(
				httphelp.GetHandleFunc(
					adminHandler.ListInventories))))

	http.HandleFunc(routing.RouteAdminAddInventory,
		httphelp.AuthenticatedHandleFunc(
			deps,
			httphelp.IsAdminHandleFunc(
				httphelp.GetPostHandleFunc(
					adminHandler.AddInventory))))

	http.HandleFunc(routing.RouteAdminDeleteInventory,
		httphelp.AuthenticatedHandleFunc(
			deps,
			httphelp.IsAdminHandleFunc(
				httphelp.PostHandleFunc(
					adminHandler.DeleteInventory))))

	http.HandleFunc(routing.RouteAdminEditInventory,
		httphelp.AuthenticatedHandleFunc(
			deps,
			httphelp.IsAdminHandleFunc(
				httphelp.GetPostHandleFunc(
					adminHandler.EditInventory))))

	log.Println("Done")
}
