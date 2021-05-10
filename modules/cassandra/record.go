package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type Record struct {
	AccountUuid gocql.UUID `cql:"account_uuid"`
	RecordUuid  gocql.UUID `cql:"record_uuid"`
	Tags        []string   `cql:"tags"`
	Words       []string   `cql:"words"`
	Record      string     `cql:"record"`
	Title       string     `cql:"title"`
	Importance  int        `cql:"importance"`
}

// UpsertRecord will insert or update a record in the Cassandra database
func (s *Session) UpsertRecord(record Record) bool {
	if err := s.Connection.Query(`
		UPDATE record 
		SET tags = tags + ?,
			words = words + ?,
			record = ?,
			title = ?,
			importance = ?
		WHERE account_uuid = ?
		  AND record_uuid = ?
		`, record.Tags, record.Words, record.Record, record.Title, record.Importance, record.AccountUuid, record.RecordUuid).Exec(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GetRecords Will return all of the records for the account in the passed in recordUuid list
func (s *Session) GetRecords(accountUuid gocql.UUID, recordUuids []gocql.UUID) []Record {
	var recordSets []Record
	m := map[string]interface{}{}
	query := "SELECT * FROM record WHERE account_uuid = ? and record_uuid IN ?"
	iterable := s.Connection.Query(query, accountUuid, recordUuids).Iter()
	for iterable.MapScan(m) {
		recordSets = append(recordSets, Record{
			AccountUuid: m["account_uuid"].(gocql.UUID),
			RecordUuid:  m["record_uuid"].(gocql.UUID),
			Tags:        m["tags"].([]string),
			Words:       m["words"].([]string),
			Record:      m["record"].(string),
			Title:       m["title"].(string),
			Importance:  m["importance"].(int),
		})
		m = map[string]interface{}{}
	}
	return recordSets
}