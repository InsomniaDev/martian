package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()

	r.HandleFunc("/record/{record}", UpdateRecord).Methods("POST")
	r.HandleFunc("/record/new", InsertNewRecord).Methods("POST")
	r.HandleFunc("/query/{record}", RetrieveRecord).Methods("GET")

	http.ListenAndServe(":8050", r)
}
