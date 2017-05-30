package main

import (
	"net/http"
)

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	app.cookieAuth.DeleteLoginCookie(w)
	http.Redirect(w, r, RouteLogin, http.StatusSeeOther)
}
