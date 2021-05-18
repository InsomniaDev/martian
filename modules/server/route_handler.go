package server

import (
	"encoding/json"
	"fmt"
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
	AccountUuid gocql.UUID
}

type MartianResponse struct {
	Records []cassandra.Record `json:"records"`
	Message string             `json:"message"`
}

// DecipherQuery point of this is decipher the query string
func DecipherQuery(w http.ResponseWriter, r *http.Request) {

	recordData, err := getMessageBody(r)
	if err != nil {
		w.Write([]byte("Not Authorized"))
		return
	}

	logic.ParseEntry(recordData.Record)

	helpCommand := strings.TrimSpace(recordData.Record)
	if strings.ToLower(helpCommand) == "help" {
		returnHelp(w, r)
	}
	// TODO: Check what the prefix of this is
}

// insertNewRecord will parse the http POST and insert a new record into the cassandra Record table
//  return success boolean
func insertNewRecord(w http.ResponseWriter, r *http.Request, message MartianBody) {
	log.Println("we are inserting new record")

	// Parse the record into a Cassandra record and then set the AccountUuid
	record := logic.ParseRecordIntoCassandraRecord(message.Record)
	record.AccountUuid = message.AccountUuid

	fmt.Printf("%#v", record)
	return

	// Insert the provided record into the Cassandra database
	inserted := logic.UpsertRecord(&CassandraConnection, record)
	if inserted {
		log.Println("Inserted new record")
	} else {
		log.Panic("Record insert failed")
	}
}

// updateRecord will go through and update the record that is provided]]]
func updateRecord(w http.ResponseWriter, r *http.Request, message MartianBody) {
	log.Println("we are updating a record")
	vars := mux.Vars(r)
	recordUuid := vars["recordUuid"]

	// Parse the record into a Cassandra record and then set the AccountUuid
	record := logic.ParseRecordIntoCassandraRecord(message.Record)
	record.AccountUuid = message.AccountUuid
	
	uuid, err := gocql.ParseUUID(recordUuid)
	if err != nil {
		log.Fatal(err)
	}
	record.RecordUuid = uuid

	// TODO: Need to search for the record here and update with the data that already exists

	// Insert the provided record into the Cassandra database
	inserted := logic.UpsertRecord(&CassandraConnection, record)
	if inserted {
		log.Println("Inserted new record")
	} else {
		log.Panic("Record insert failed")
	}
}

// returnHelp will return the help output which specifies how to use the application
func returnHelp(w http.ResponseWriter, r *http.Request) {
	helpOutput := `
		To start:
			A query begins with one of the W's and an H -> the Who, What, When, Where, Why, How
			NEW Record begins with anything else
			UPDATE begins with the record Uuid and then it will fully replace that record. 
			DELETE with the record UUID will delte the record
	`
	var response MartianResponse
	response.Message = helpOutput

	// Convert response to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// retrieveRecord will retrieve the records matching the query string and return them to the calling application
func retrieveRecord(w http.ResponseWriter, r *http.Request, message MartianBody) {

	// Set the account uuid
	var searchRecord logic.RecordRequest
	searchRecord.AccountUuid = message.AccountUuid

	// Parse out the tags and words from the passed record
	entities, words := logic.ParseEntry(message.Record)
	searchRecord.Entities = entities
	searchRecord.Words = words

	// Retrieve the records that match the incoming request
	searchRecord.ParseRequest(&CassandraConnection, &CommonWords, 3)
	data := searchRecord.RetrieveRecords(&CassandraConnection, 3)

	var response MartianResponse
	response.Message = "Consumed: " + message.Record
	response.Records = data

	// Convert response to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// getAuthToken parses the header and returns the auth token string
func getAuthToken(r *http.Request) (gocql.UUID, error) {
	authToken, err := gocql.ParseUUID(r.Header.Get("x-access-token"))
	if err != nil {
		return authToken, err
	}
	return authToken, nil
}

// getMessageBody parses the http request body into a MartianBody struct and returns the struct and the auth token string
func getMessageBody(r *http.Request) (MartianBody, error) {
	var body MartianBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
	}
	body.AccountUuid, err = getAuthToken(r)

	return body, err
}
