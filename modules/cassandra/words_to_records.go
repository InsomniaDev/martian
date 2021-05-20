package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type WordsToRecord struct {
	Word        string     `cql:"word"`
	AccountUuid gocql.UUID `cql:"account_uuid"`
	RecordUuid  string     `cql:"record_uuid"`
}

type WordsToRecords struct {
	Word        string     `cql:"word"`
	AccountUuid gocql.UUID `cql:"account_uuid"`
	RecordUuid  []string   `cql:"record_uuid"`
}

// UpsertWordsToRecord will add record association to word
func (s *Session) UpsertWordsToRecord(words WordsToRecord) {
	if err := s.Connection.Query(`
		UPDATE words_to_records
		SET record_uuid = record_uuid + ?
		WHERE account_uuid = ?
		  AND word = ?
		`, words.RecordUuid, words.AccountUuid, words.Word).Exec(); err != nil {
		fmt.Println(err)
	}
}

// UpsertMultipleWordsToRecord will upsert multiple words and records into the tables
func (s *Session) UpsertMultipleWordsToRecord(words []WordsToRecord) {
	batchQuery := s.Connection.NewBatch(gocql.LoggedBatch)
	for _, wordRecord := range words {
		recordUuid := []string{wordRecord.RecordUuid}
		batchQuery.Query(`
		UPDATE words_to_records
		SET record_uuid = record_uuid + ?
		WHERE account_uuid = ?
		  AND word = ?`, recordUuid, wordRecord.AccountUuid, wordRecord.Word)
	}
	if err := s.Connection.ExecuteBatch(batchQuery); err != nil {
		fmt.Println(err)
	}
}

// DeleteRecordsFromWords will delete the associated record_uuid from the entry
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

// GetWordsToRecords will get all matching words for the associated account
func (s *Session) GetWordsToRecords(words []string, account gocql.UUID) []WordsToRecords {
	var wordsToRecords []WordsToRecords
	m := map[string]interface{}{}
	query := "SELECT * FROM words_to_records WHERE account_uuid = ? and word IN ?"
	fmt.Println("\n\nQUERY:", query, account, words)
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
