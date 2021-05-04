package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type MartianBody struct {
	RecordUuid  string `json:"recordUuid"`
	MessageUuid string `json:"messageUuid"`
	Record      string `json:"entry"`
}

func InsertNewRecord(w http.ResponseWriter, r *http.Request) {
	log.Panic("not implemented")
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	record := vars["record"]
	log.Print(record)
	log.Panic("not implemented")
}

func RetrieveRecord(w http.ResponseWriter, r *http.Request) {

	var body MartianBody
	authToken := r.Header.Get("x-access-token")

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(authToken)
}
