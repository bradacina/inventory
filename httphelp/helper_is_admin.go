package httphelp

import (
	"net/http"
)

func IsAdminHandler(handler http.Handler) http.HandlerFunc {
	return makeIsAdminHandler(handler.ServeHTTP)
}

func IsAdminHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeIsAdminHandler(handlerFunc)
}

func makeIsAdminHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetLoggedInUser(r)
		if !user.IsAdmin {
			StatusCode(w, http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
