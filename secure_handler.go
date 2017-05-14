package main

import (
	"net/http"
)

func (app *app) secure(w http.ResponseWriter, r *http.Request) {
	user := GetLoggedInUser(r)

	w.Write([]byte("Welcome " + user.Email))
}
