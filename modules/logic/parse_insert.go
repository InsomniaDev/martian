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
	conn.UpsertMultipleWordsToRecord(wordRecords)

	// TODO: need to upsert the entities to the record here

	return inserted
}
