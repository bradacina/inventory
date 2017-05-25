package main

import (
	"net/http"

	"github.com/bradacina/inventory/httphelp"
)

func (app *app) secure(w http.ResponseWriter, r *http.Request) {
	user := httphelp.GetLoggedInUser(r)

	w.Write([]byte("Welcome " + user.Email))
}
