package logic

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
	"github.com/jdkato/prose/v2"
)

func SortAndRetrieveRecordUuids(records [][]string, numOfRecordsToReturn int) []string {

	// Go through all of the passed in records and create counts of times the records appear
	recordCounts := make(map[string]int)
	for i := range records {
		for _, recordUUID := range records[i] {
			if _, ok := recordCounts[recordUUID]; ok {
				recordCounts[recordUUID] = recordCounts[recordUUID] + 1
			} else {
				recordCounts[recordUUID] = 1
			}
		}
	}

	// Sort all of the values descending so that we get the top found elements first
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range recordCounts {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// Create a list of the top records that exists of the passed in num
	recordUuidsToReturn := []string{}
	lengthOfRecords := len(ss)
	for i := 0; i < numOfRecordsToReturn && i < lengthOfRecords; i++ {
		recordUuidsToReturn = append(recordUuidsToReturn, ss[i].Key)
	}

	return recordUuidsToReturn
}

func ParseRecordIntoCassandraRecord(postRecord string) cassandra.Record {
	var record cassandra.Record
	randomUuid, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}
	record.RecordUuid = randomUuid

	// Set the title as the first line
	record.Title = strings.Split(postRecord, "\n")[0]

	// Set the record to the whole entry including the title
	record.Record = postRecord

	fmt.Println("\nrecord to analyze:", postRecord)

	// Parse out the tags and words from the passed record
	entities, words := ParseEntry(postRecord)
	record.Entities = entities
	record.Words = words

	// Set importance to 0 since this is the first insert
	record.Importance = 0
	return record
}

// ParseEntry Split up the incoming query record between words and tags
func ParseEntry(recordData string) ([]string, []string) {
	// Symbols required to move
	symbolsToRemove := []string{"~", "`", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "<", ">", "?", "/", ";", ":", "{", "}", "[", "]", "|", "=", "+", "\n", "-"}

	// Remove all symbols not needed first
	for _, symbol := range symbolsToRemove {
		recordData = strings.ReplaceAll(recordData, symbol, "")
	}
	
	// Make sure that we are always checking against lowercase entries
	// recordData = strings.ToLower(recordData)

	// Get the language from the document
	// https://github.com/jdkato/prose
	doc, err := prose.NewDocument(recordData)
	if err != nil {
		log.Println("ERROR", err)
		return nil, nil
	}

	fmt.Println("\n\ntokens/words", doc.Tokens())
	fmt.Println("\n\nentities", doc.Entities())

	fmt.Println("\n\nsentences",doc.Sentences())

	var tags []string
	var words []string
	for _, token := range doc.Tokens() {
		words = append(words, strings.ToLower(token.Text))
	}
	for _, entity := range doc.Entities() {
		tags = append(tags, entity.Text)
	}
	return tags, words
}
