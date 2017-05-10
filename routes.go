package main

import (
	"log"
	"net/http"
)

func registerRoutes(app *app) {
	log.Println("Setting up routes...")
	http.HandleFunc("/register", app.register)
	http.HandleFunc("/login", app.login)
	log.Println("Done")
}
