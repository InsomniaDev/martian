package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/insomniadev/martian/modules/cassandra"
	"github.com/insomniadev/martian/modules/logic"
)

type MartianBody struct {
	RecordUuid  string `json:"recordUuid"`
	MessageUuid string `json:"messageUuid"`
	Record      string `json:"entry"`
}

// InsertNewRecord will insert a new record into the cassandra Record table and will return success boolean
func InsertNewRecord(w http.ResponseWriter, r *http.Request) {
	recordData, authToken := getMessageBody(r)
	accountUuid, err := gocql.ParseUUID(authToken)
	if err != nil {
		log.Fatal(err)
	}

	// Create the record and assign the account UUID and create a new record UUID
	var record cassandra.Record
	record.AccountUuid = accountUuid
	record.RecordUuid, err = gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}

	// Set the title as the first line
	record.Title = strings.Split(recordData.Record, "\n")[0]

	// Set the record to the whole entry including the title
	record.Record = recordData.Record

	// Parse out the tags and words from the passed record
	tags, words := parseEntry(recordData.Record)
	record.Tags = tags
	record.Words = words

	// Set importance to 0 since this is the first insert
	record.Importance = 0

	// Insert the provided record into the Cassandra database
	inserted := logic.UpsertRecord(&CassandraConnection, record)
	if inserted {
		log.Println("Inserted new record")
	} else {
		log.Panic("Record insert failed")
	}
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	record := vars["recordUuid"]
	log.Print(record)
	log.Panic("not implemented")
}

// RetrieveRecord will retrieve the records matching the query string and return them to the calling application
func RetrieveRecord(w http.ResponseWriter, r *http.Request) {
	recordData, authToken := getMessageBody(r)
	accountUuid, err := gocql.ParseUUID(authToken)
	if err != nil {
		log.Fatal(err)
	}

	// Set the account uuid
	var searchRecord logic.RecordRequest
	searchRecord.AccountUuid = accountUuid

	// Parse out the tags and words from the passed record
	tags, words := parseEntry(recordData.Record)
	searchRecord.Tags = tags
	searchRecord.Words = words

	// Retrieve the records that match the incoming request
	searchRecord.ParseRequest(&CassandraConnection, 3)
	data := searchRecord.RetrieveRecords(&CassandraConnection, 3)

	// Convert response to JSON
	jsonData, err := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getAuthToken(r *http.Request) string {
	return r.Header.Get("x-access-token")
}

func getMessageBody(r *http.Request) (MartianBody, string) {
	var body MartianBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}
	return body, getAuthToken(r)
}

// Split up the incoming query record between words and tags
func parseEntry(recordData string) ([]string, []string) {
	// split the string into an array first
	recordDataSlice := strings.Split(recordData, " ")

	var tags []string
	var words []string
	// Take apart and get separate lists of tags and words
	for _, value := range recordDataSlice {
		if strings.HasPrefix(value, "#") {
			tags = append(tags, value)
		} else {
			words = append(words, value)
		}
	}
	return tags, words
}
