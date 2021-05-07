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

	r.HandleFunc("/record/update/{recordUuid}", UpdateRecord).Methods("POST")
	r.HandleFunc("/record/new", InsertNewRecord).Methods("POST")
	r.HandleFunc("/query", RetrieveRecord).Methods("POST")

	http.ListenAndServe(":9000", r)
}
