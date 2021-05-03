package logic

import (
	"sort"
)

func SortAndRetrieveRecordUuids(records [][]string, numOfRecordsToReturn int) []string {

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

	recordUuidsToReturn := []string{}
	for i := 0; i < numOfRecordsToReturn; i++ {
		recordUuidsToReturn = append(recordUuidsToReturn, ss[i].Key)
	}

	return recordUuidsToReturn
}
