package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type TagsToRecords struct {
	Tag         string
	AccountUuid gocql.UUID
	RecordUuid  []string
}

func (s *Session) UpsertTagsToRecords(tags TagsToRecords) {
	if err := s.Connection.Query(`
		UPDATE tags_to_records 
		SET record_uuid = record_uuid + ?
		WHERE account_uuid = ?
		  AND tag = ?
		`, tags.RecordUuid, tags.AccountUuid, tags.Tag).Exec(); err != nil {
		fmt.Println(err)
	}
}

func (s *Session) DeleteRecordsFromTags(tags TagsToRecords) {
	if err := s.Connection.Query(`
		UPDATE tags_to_records 
		SET record_uuid = record_uuid - ?
		WHERE account_uuid = ?
		  AND tag = ?
		`, tags.RecordUuid, tags.AccountUuid, tags.Tag).Exec(); err != nil {
		fmt.Println(err)
	}
}
