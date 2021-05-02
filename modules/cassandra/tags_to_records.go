package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type TagsToRecords struct {
	Tag         string     `cql:"tag"`
	AccountUuid gocql.UUID `cql:"account_uuid"`
	RecordUuid  []string   `cql:"record_uuid"`
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

func (s *Session) GetTagsToRecords(tags []string, account gocql.UUID) ([]TagsToRecords) {
	var tagsToRecords []TagsToRecords
	m := map[string]interface{}{}
	query := "SELECT * FROM tags_to_records WHERE account_uuid = ? and tag IN ?"
	iterable := s.Connection.Query(query, account, tags).Iter()
	for iterable.MapScan(m) {
		tagsToRecords = append(tagsToRecords, TagsToRecords{
			Tag:        m["tag"].(string),
			AccountUuid: m["account_uuid"].(gocql.UUID),
			RecordUuid:  m["record_uuid"].([]string),
		})
		m = map[string]interface{}{}
	}
	return tagsToRecords
}
