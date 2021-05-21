package logic

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

// TODO: Process an incoming request

func ParseInsert(conn *cassandra.Session, accountUuid gocql.UUID) {
	fmt.Println("I don't do anything yet")

}

func UpsertRecord(conn *cassandra.Session, record cassandra.Record) bool {

	inserted := conn.UpsertRecord(record)

	// Create the word association with the record
	var wordRecords []cassandra.WordsToRecord
	for _, word := range record.Words {
		wordRecords = append(wordRecords, cassandra.WordsToRecord{Word: word, AccountUuid: record.AccountUuid, RecordUuid: record.RecordUuid.String()})
	}
	batchQuery := conn.UpsertMultipleWordsToRecord(wordRecords)

	// Create the entity association with the record
	var entityRecords []cassandra.EntitiesToRecord
	for _, entity := range record.Entities {
		entityRecords = append(entityRecords, cassandra.EntitiesToRecord{Entity: entity, AccountUuid: record.AccountUuid, RecordUuid: record.RecordUuid.String()})
	}
	batchQuery = conn.UpsertMultipleEntitiesToRecord(entityRecords, batchQuery)

	conn.ExecuteBatch(batchQuery)

	// TODO: need to upsert the entities to the record here

	return inserted
}
