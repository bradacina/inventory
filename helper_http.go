package main

import (
	"net/http"
)

func StatusCode(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}
