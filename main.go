package main

import (
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":9080", nil))
}
