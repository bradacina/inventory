package main

import (
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		serveTemplate(w, TemplateRegister)

	} else if r.Method == http.MethodPost {

	} else {
		http.NotFound(w, r)
	}
}
