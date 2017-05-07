package main

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		serveTemplate(w, TemplateLogin)
	} else if r.Method == http.MethodPost {

	} else {
		http.NotFound(w, r)
	}

}
