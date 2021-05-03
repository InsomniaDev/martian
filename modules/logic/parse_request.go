package logic

import (
	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

type RecordRequest struct {
	AccountUuid gocql.UUID
	Tags        []string
	Words       []string
	Records     []string
}

// ParseRequest gets all of the records that match the tags and words and sets the Records variable
func (rr *RecordRequest) ParseRequest(conn *cassandra.Session, numOfRecords int) {
	var likelyRecords []string

	// Get three records from the provided words
	recordsFromWords := RetrieveListOfRecordsForWords(conn, rr.AccountUuid, rr.Words, 3)

	if len(rr.Tags) > 0 {
		// If tags are provided then get three records from the provided tags
		recordsFromTags := RetrieveListOfRecordsForTags(conn, rr.AccountUuid, rr.Tags, 3)
		likelyRecords = returnMostImportantRecords(recordsFromWords, recordsFromTags)
	} else {
		likelyRecords = returnMostImportantRecords(recordsFromWords, nil)
	}

	var recordsToReturn []string
	if len(likelyRecords) > 0 {
		for i := 0; i < numOfRecords && i < len(recordsToReturn); i++ {
			recordsToReturn = append(recordsToReturn, likelyRecords[i])
		}
	}
	rr.Records = recordsToReturn
}

// // RetrieveRecords
// func (rr *RecordRequest)

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
