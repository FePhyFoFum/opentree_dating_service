package main

import (
	"log"
	"net/http"

	"github.com/FePhyFoFum/opentree_dating_service"
)

func main() {

	router := inducedates.NewRouter()

	log.Fatal(http.ListenAndServe(":10999", router))
}
