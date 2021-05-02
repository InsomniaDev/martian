package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type WordsToRecords struct {
	Word        string     `cql:"word"`
	AccountUuid gocql.UUID `cql:"account_uuid"`
	RecordUuid  []string   `cql:"record_uuid"`
}

func (s *Session) UpsertWordsToRecords(words WordsToRecords) {
	if err := s.Connection.Query(`
		UPDATE words_to_records
		SET record_uuid = record_uuid + ?
		WHERE account_uuid = ?
		  AND word = ?
		`, words.RecordUuid, words.AccountUuid, words.Word).Exec(); err != nil {
		fmt.Println(err)
	}
}

func (s *Session) DeleteRecordsFromWords(words WordsToRecords) {
	if err := s.Connection.Query(`
		UPDATE words_to_records 
		SET record_uuid = record_uuid - ?
		WHERE account_uuid = ?
		  AND word = ?
		`, words.RecordUuid, words.AccountUuid, words.Word).Exec(); err != nil {
		fmt.Println(err)
	}
}

func (s *Session) GetWordsToRecords(words []string, account gocql.UUID) ([]WordsToRecords) {
	var wordsToRecords []WordsToRecords
	m := map[string]interface{}{}
	query := "SELECT * FROM words_to_records WHERE account_uuid = ? and word IN ?"
	iterable := s.Connection.Query(query, account, words).Iter()
	for iterable.MapScan(m) {
		wordsToRecords = append(wordsToRecords, WordsToRecords{
			Word:        m["word"].(string),
			AccountUuid: m["account_uuid"].(gocql.UUID),
			RecordUuid:  m["record_uuid"].([]string),
		})
		m = map[string]interface{}{}
	}
	return wordsToRecords
}
