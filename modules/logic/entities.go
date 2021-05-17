package logic

import (
	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

// RetrieveListOfRecordsForEntities will get the defined number of relevant records by the entities provided and return the record uuids
func RetrieveListOfRecordsForEntities(conn *cassandra.Session, accountUuid gocql.UUID, entities []string, numOfRecords int) []string {

	// Get all of the results from the Cassandra database that have the entities
	results := conn.GetEntitiesToRecords(entities, accountUuid)

	// Filter down to just the records that were retrieved
	var recordsRetrieved [][]string
	for _, records := range results {
		recordsRetrieved = append(recordsRetrieved, records.RecordUuid)
	}

	// Get the top list of received records
	return SortAndRetrieveRecordUuids(recordsRetrieved, numOfRecords)
}
