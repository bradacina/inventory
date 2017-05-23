package main

import (
	"context"
	"log"
	"net/http"
)

type loggedInUserKeyType string

const loggedInUserKey = loggedInUserKeyType("loggedInUserKey")

func (app *app) authenticatedHandler(handler http.Handler) http.HandlerFunc {
	return makeAuthHandler(app, handler.ServeHTTP)
}

func (app *app) authenticatedHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeAuthHandler(app, handlerFunc)
}

func makeAuthHandler(app *app, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loginInfo, err := app.cookieHelper.getLoginCookie(r)
		if err != nil {
			log.Println("Login cookie not present. Redirecting to login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user, err := app.userService.GetByEmail(loginInfo.Username)
		if err != nil {
			log.Println("Login cookie contained non existant username")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(context.Background(), loggedInUserKey, user)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

func GetLoggedInUser(r *http.Request) *User {
	user := r.Context().Value(loggedInUserKey).(*User)
	return user
}
