package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type WordsToRecords struct {
	Word        string
	AccountUuid gocql.UUID
	RecordUuid  []string
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
