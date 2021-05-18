package server

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/insomniadev/martian/modules/cassandra"
)

var CassandraConnection cassandra.Session
var CommonWords []string

func Start() {
	r := mux.NewRouter()

	CassandraConnection = cassandra.Session{}
	CassandraConnection.Init()
	CommonWords = strings.Split(CassandraConnection.GetConfig("commonWords"), ",")
	// Setup variables for efficiency, do once and use everywhere
	sort.Strings(CommonWords)

	// r.HandleFunc("/record/update/{recordUuid}", UpdateRecord).Methods("POST")
	// r.HandleFunc("/record/new", InsertNewRecord).Methods("POST")
	r.HandleFunc("/query", DecipherQuery).Methods("POST")

	http.ListenAndServe(":9000", r)
}
