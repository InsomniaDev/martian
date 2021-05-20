package logic

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

type RecordRequest struct {
	AccountUuid gocql.UUID
	Entities        []string
	Words       []string
	Records     []gocql.UUID
}

// ParseRequest gets all of the records that match the entities and words and sets the Records variable
func (rr *RecordRequest) ParseRequest(conn *cassandra.Session, commonWords *[]string, numOfRecords int) {
	var likelyRecords []string

	// Get three records from the provided words
	recordsFromWords := RetrieveListOfRecordsForWords(conn, rr.AccountUuid, rr.Words, commonWords, 3)
	fmt.Println("\n\nrecordsFromWords",recordsFromWords)

	if len(rr.Entities) > 0 {
		// If entities are provided then get three records from the provided entities
		recordsFromEntities := RetrieveListOfRecordsForEntities(conn, rr.AccountUuid, rr.Entities, 3)
		likelyRecords = returnMostImportantRecords(recordsFromWords, recordsFromEntities)
	} else {
		likelyRecords = returnMostImportantRecords(recordsFromWords, nil)
	}

	var recordsToReturn []gocql.UUID
	if len(likelyRecords) > 0 {
		for i := 0; i < numOfRecords && i < len(recordsToReturn); i++ {
			uuid, err := gocql.ParseUUID(likelyRecords[i])
			if err != nil {
				fmt.Println(err)
			}
			recordsToReturn = append(recordsToReturn, uuid)
		}
	}
	rr.Records = recordsToReturn
}

// RetrieveRecords will retrieve the records in the RecordRequest object, numOfRecords will limit returned amount of records
func (rr *RecordRequest) RetrieveRecords(conn *cassandra.Session, numOfRecords int) []cassandra.Record {
	var recordsRequired []gocql.UUID
	recordLength := len(rr.Records)
	for i := 0; i < numOfRecords && i < recordLength; i++ {
		recordsRequired = append(recordsRequired, rr.Records[i])
	}

	return conn.GetRecords(rr.AccountUuid, recordsRequired)
}

/* returnMostImportantRecords goes through and returns the most important records according to:
*		word and tag records match
*		tag records are next
*		word records are next
 */
func returnMostImportantRecords(wordRecords []string, tagRecords []string) []string {

	words := wordRecords
	tags := tagRecords

	// Check against if there are matches between tag records and word records
	var likelyRecordUuids []string
	if tagRecords != nil {
		words = []string{}
		for _, wordRecord := range wordRecords {
			found := false
			for _, tagRecord := range tagRecords {
				if wordRecord == tagRecord {
					likelyRecordUuids = append(likelyRecordUuids, wordRecord)
					found = true
				}
			}
			if !found {
				words = append(words, wordRecord)
			}
		}
	}

	// If there were matches then we need to remove those matches from the tags
	if len(likelyRecordUuids) > 0 {
		tags = []string{}
		for _, tagRecord := range tagRecords {
			found := false
			for _, likelyRecord := range likelyRecordUuids {
				if tagRecord == likelyRecord {
					found = true
				}
			}
			if !found {
				tags = append(tags, tagRecord)
			}
		}
	}

	// Add together all of the lists into one and return it, adding tags first
	if tagRecords != nil {
		likelyRecordUuids = append(likelyRecordUuids, tags...)
	}

	// Add the top word records to the list to return
	likelyRecordUuids = append(likelyRecordUuids, words...)

	return likelyRecordUuids
}
