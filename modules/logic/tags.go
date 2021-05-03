package logic

import (
	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

// RetrieveListOfRecordsForTags will get the defined number of relevant records by the tags provided and return the record uuids
func RetrieveListOfRecordsForTags(conn *cassandra.Session, accountUuid gocql.UUID, tags []string, numOfRecords int) []string {

	// Get all of the results from the Cassandra database that have the Tags
	results := conn.GetTagsToRecords(tags, accountUuid)

	// Filter down to just the records that were retrieved
	var recordsRetrieved [][]string
	for _, records := range results {
		recordsRetrieved = append(recordsRetrieved, records.RecordUuid)
	}

	// Get the top list of received records
	return SortAndRetrieveRecordUuids(recordsRetrieved, numOfRecords)
}
