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

// InsertNewRecord will parse the http POST and insert a new record into the cassandra Record table
//  return success boolean
func insertNewRecord(w http.ResponseWriter, r *http.Request, message MartianBody) {
	log.Println("we are inserting new record")

	// Create the record and assign the account UUID and create a new record UUID
	var record cassandra.Record
	record.AccountUuid = message.AccountUuid
	randomUuid, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}
	record.RecordUuid = randomUuid

	// Set the title as the first line
	record.Title = strings.Split(message.Record, "\n")[0]

	// Set the record to the whole entry including the title
	record.Record = message.Record

	// Parse out the tags and words from the passed record
	tags, words := parseEntry(message.Record)
	record.Tags = tags
	record.Words = words

	// Set importance to 0 since this is the first insert
	record.Importance = 0
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

func updateRecord(w http.ResponseWriter, r *http.Request, message MartianBody) {
	log.Println("we are updating a record")
	vars := mux.Vars(r)
	recordUuid := vars["recordUuid"]

	// Create the record and assign the account UUID and create a new record UUID
	var record cassandra.Record
	record.AccountUuid = message.AccountUuid
	// if record.RecordUuid {
	uuid, err := gocql.ParseUUID(recordUuid)
	if err != nil {
		log.Fatal(err)
	}
	record.RecordUuid = uuid
	// }
	record.RecordUuid, err = gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}

	// Set the title as the first line
	record.Title = strings.Split(message.Record, "\n")[0]

	// Set the record to the whole entry including the title
	record.Record = message.Record

	// Parse out the tags and words from the passed record
	tags, words := parseEntry(message.Record)
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

// DecipherQuery point of this is decipher the query string
func DecipherQuery(w http.ResponseWriter, r *http.Request) {

	recordData, err := getMessageBody(r)
	if err != nil {
		w.Write([]byte("Not Authorized"))
		return
	}
	fmt.Print(recordData.Record)

	helpCommand := strings.TrimSpace(recordData.Record)
	if strings.ToLower(helpCommand) == "help" {

	}
}

// returnHelp will return the help output which specifies how to use the application
func returnHelp(w http.ResponseWriter, r *http.Request) {
	helpOutput := `
		To start:
			A query begins with zzzz
			Record Begins with xxxxx
			Update begins with zzzzz
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
	tags, words := parseEntry(message.Record)
	searchRecord.Tags = tags
	searchRecord.Words = words

	// Retrieve the records that match the incoming request
	searchRecord.ParseRequest(&CassandraConnection, 3)
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
