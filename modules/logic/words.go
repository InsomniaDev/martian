package logic

import (
	"fmt"
	"sort"

	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

// var commonWords []string

// RetrieveListOfRecordsForWords will get the defined number of relevant records and return the record uuids
func RetrieveListOfRecordsForWords(conn *cassandra.Session, accountUuid gocql.UUID, searchString []string, commonWords *[]string, numOfRecords int) []string {

	// Get the list of common words from the database
	// commonWords = commonWords

	// Remove all of the common words for the words that will be searched by
	wordsToSearch := removeCommonWords(searchString, commonWords)
	fmt.Println("\n\nwordsToSearch", wordsToSearch)

	// Get all of the results from the Cassandra database that have the words
	results := conn.GetWordsToRecords(wordsToSearch, accountUuid)

	// Filter down to just the records that were retrieved
	var recordsRetrieved [][]string
	for _, records := range results {
		recordsRetrieved = append(recordsRetrieved, records.RecordUuid)
	}

	// Get the top list of received records
	return SortAndRetrieveRecordUuids(recordsRetrieved, numOfRecords)
}

// removeCommonWords will remove the common words from the string of words that was provided
func removeCommonWords(stringToParse []string, commonWords *[]string) []string {

	// Get a list of the words that we can now search by
	goodWords := []string{}
	for _, a := range stringToParse {
		// Only if the word is not common will we add it to the list of words to check
		if !checkIfWordIsCommon(a, commonWords) {
			// Add to list
			goodWords = append(goodWords, a)
		}
	}
	return goodWords
}

// checkIfWordIsCommon does a binary search and returns if it is a common word
func checkIfWordIsCommon(word string, commonWords *[]string) bool {
	i := sort.Search(len(*commonWords),
		func(i int) bool { return (*commonWords)[i] >= word })
	if i < len(*commonWords) && (*commonWords)[i] == word {
		return true
	}
	return false
}
