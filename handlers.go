package inducedates

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//Index wasted index function
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

//InduceSubtree get induced subtree
//curl -X POST http://localhost:8080/induce_subtree -H "content-type:application/json" -d '{"ott_ids":[292466, 267845, 666104, 316878, 102710]}'
func InduceSubtree(w http.ResponseWriter, r *http.Request) {
	var query Query
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &query); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	t := GetInducedTree(query.Ottids)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
