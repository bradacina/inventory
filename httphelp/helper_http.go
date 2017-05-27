package httphelp

import (
	"net/http"
)

func StatusCode(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func PostHandler(handler http.Handler) http.HandlerFunc {
	return makeHandler(handler.ServeHTTP, http.MethodPost)
}

func PostHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeHandler(handlerFunc, http.MethodPost)
}

func GetHandler(handler http.Handler) http.HandlerFunc {
	return makeHandler(handler.ServeHTTP, http.MethodGet)
}

func GetHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeHandler(handlerFunc, http.MethodGet)
}

func GetPostHandler(handler http.Handler) http.HandlerFunc {
	return makeHandler(handler.ServeHTTP, http.MethodGet, http.MethodPost)
}

func GetPostHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeHandler(handlerFunc, http.MethodGet, http.MethodPost)
}

func makeHandler(next http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		found := false
		for i := range methods {
			if r.Method == methods[i] {
				found = true
				break

			}
		}

		if !found {
			StatusCode(w, http.StatusNotFound)
			return
		}

		next(w, r)
	}

}
