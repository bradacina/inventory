package main

import (
	"log"
	"net/http"
)

func init() {
	log.Println("Setting up routes...")
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	log.Println("Done")
}
