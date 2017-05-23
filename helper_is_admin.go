package main

import (
	"net/http"
)

func (app *app) isAdminHandler(handler http.Handler) http.HandlerFunc {
	return makeIsAdminHandler(app, handler.ServeHTTP)
}

func (app *app) isAdminHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeIsAdminHandler(app, handlerFunc)
}

func makeIsAdminHandler(app *app, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetLoggedInUser(r)
		if !user.IsAdmin {
			StatusCode(w, http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
