package main

import (
	"inducedates"
	"log"
	"net/http"
)

func main() {

	router := inducedates.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
