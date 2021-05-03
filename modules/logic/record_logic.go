package logic

import (
	"sort"
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
