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

func (s *Session) UpsertRecord(record Record) {
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
	}
}

// func (s *Session) DeleteRecord(tags Record) {
// 	if err := s.Connection.Query(`
// 		UPDATE tags_to_records
// 		SET record_uuid = record_uuid - ?
// 		WHERE account_uuid = ?
// 		  AND tag = ?
// 		`, tags.RecordUuid, tags.AccountUuid, tags.Tag).Exec(); err != nil {
// 		fmt.Println(err)
// 	}
// }
