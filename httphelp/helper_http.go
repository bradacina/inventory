package httphelp

import (
	"net/http"
)

func StatusCode(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func PostHandler(handler http.Handler) http.HandlerFunc {
	return makePostHandler(handler.ServeHTTP)
}

func PostHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makePostHandler(handlerFunc)
}

func GetHandler(handler http.Handler) http.HandlerFunc {
	return makeGetHandler(handler.ServeHTTP)
}

func GetHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeGetHandler(handlerFunc)
}

func makePostHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			StatusCode(w, http.StatusNotFound)
			return
		}

		next(w, r)
	}
}

func makeGetHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			StatusCode(w, http.StatusNotFound)
			return
		}

		next(w, r)
	}
}
