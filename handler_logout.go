package main

import (
	"net/http"
)

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	app.cookieHelper.deleteLoginCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
