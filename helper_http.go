package main

import (
	"net/http"
)

func StatusCode(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *app) postHandler(handler http.Handler) http.HandlerFunc {
	return makePostHandler(app, handler.ServeHTTP)
}

func (app *app) postHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makePostHandler(app, handlerFunc)
}

func (app *app) getHandler(handler http.Handler) http.HandlerFunc {
	return makeGetHandler(app, handler.ServeHTTP)
}

func (app *app) getHandleFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return makeGetHandler(app, handlerFunc)
}

func makePostHandler(app *app, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			StatusCode(w, http.StatusNotFound)
			return
		}

		next(w, r)
	}
}

func makeGetHandler(app *app, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			StatusCode(w, http.StatusNotFound)
			return
		}

		next(w, r)
	}
}
