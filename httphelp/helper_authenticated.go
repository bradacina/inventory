package httphelp

import (
	"context"
	"log"
	"net/http"

	"github.com/bradacina/inventory/db"
	"github.com/bradacina/inventory/logincookie"
	"github.com/bradacina/inventory/services"
)

type loggedInUserKeyType string

const loggedInUserKey = loggedInUserKeyType("loggedInUserKey")

type AuthDependencyInjector interface {
	GetCookieAuthentication() *logincookie.CookieAuthentication
	GetUserService() services.UserServicer
}

func AuthenticatedHandler(authDI AuthDependencyInjector, handler http.Handler) http.HandlerFunc {
	return makeAuthHandler(authDI, handler.ServeHTTP)
}

func AuthenticatedHandleFunc(authDI AuthDependencyInjector, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeAuthHandler(authDI, handlerFunc)
}

func makeAuthHandler(authDI AuthDependencyInjector,
	handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		loginInfo, err := authDI.GetCookieAuthentication().GetLoginCookie(r)
		if err != nil {
			log.Println("Login cookie not present. Redirecting to login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user, err := authDI.GetUserService().GetByEmail(loginInfo.Username)
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

func GetLoggedInUser(r *http.Request) *db.User {
	user := r.Context().Value(loggedInUserKey).(*db.User)
	return user
}
