package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type EntitiesToRecord struct {
	Entity      string     `cql:"entity"`
	AccountUuid gocql.UUID `cql:"account_uuid"`
	RecordUuid  string   `cql:"record_uuid"`
}

type EntitiesToRecords struct {
	Entity      string     `cql:"entity"`
	AccountUuid gocql.UUID `cql:"account_uuid"`
	RecordUuid  []string   `cql:"record_uuid"`
}

// UpsertEntitiesToRecords will add record association to tag
func (s *Session) UpsertEntitiesToRecords(tags EntitiesToRecords) {
	if err := s.Connection.Query(`
		UPDATE entities_to_records 
		SET record_uuid = record_uuid + ?
		WHERE account_uuid = ?
		  AND entity = ?
		`, tags.RecordUuid, tags.AccountUuid, tags.Entity).Exec(); err != nil {
		fmt.Println(err)
	}
}

// UpsertMultipleWordsToRecord will upsert multiple words and records into the tables
func (s *Session) UpsertMultipleEntitiesToRecord(words []EntitiesToRecord, batchQuery *gocql.Batch) *gocql.Batch {
	for _, entityRecord := range words {
		recordUuid := []string{entityRecord.RecordUuid}
		batchQuery.Query(`
		UPDATE entities_to_records
		SET record_uuid = record_uuid + ?
		WHERE account_uuid = ?
		  AND entity = ?`, recordUuid, entityRecord.AccountUuid, entityRecord.Entity)
	}
	return batchQuery
}

// DeleteRecordsFromEntities will delete the records from the Entities
func (s *Session) DeleteRecordsFromEntities(tags EntitiesToRecords) {
	if err := s.Connection.Query(`
		UPDATE entities_to_records 
		SET record_uuid = record_uuid - ?
		WHERE account_uuid = ?
		  AND entity = ?
		`, tags.RecordUuid, tags.AccountUuid, tags.Entity).Exec(); err != nil {
		fmt.Println(err)
	}
}

// Will get all records that have the provided tags
func (s *Session) GetEntitiesToRecords(tags []string, account gocql.UUID) []EntitiesToRecords {
	var entitiesToRecords []EntitiesToRecords
	m := map[string]interface{}{}
	query := "SELECT * FROM entities_to_records WHERE account_uuid = ? and tag IN ?"
	iterable := s.Connection.Query(query, account, tags).Iter()
	for iterable.MapScan(m) {
		entitiesToRecords = append(entitiesToRecords, EntitiesToRecords{
			Entity:      m["entity"].(string),
			AccountUuid: m["account_uuid"].(gocql.UUID),
			RecordUuid:  m["record_uuid"].([]string),
		})
		m = map[string]interface{}{}
	}
	return entitiesToRecords
}
