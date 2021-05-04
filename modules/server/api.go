package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insomniadev/martian/modules/cassandra"
)

var CassandraConnection cassandra.Session

func Start() {
	r := mux.NewRouter()

	CassandraConnection = cassandra.Session{}
	CassandraConnection.Init()

	r.HandleFunc("/record/{record}", UpdateRecord).Methods("POST")
	r.HandleFunc("/record/new", InsertNewRecord).Methods("POST")
	r.HandleFunc("/query/{record}", RetrieveRecord).Methods("GET")

	http.ListenAndServe(":8050", r)
}
