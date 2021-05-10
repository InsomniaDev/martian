package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cache"
	"github.com/insomniadev/martian/modules/cassandra"
	"github.com/insomniadev/martian/modules/server"
)

func main() {
	// testLocalCache()
	server.Start()
}

func testLocalCache() {
	localCache := cache.LocalCache{}
	localCache.Init()

	type something struct {
		value string
	}
	testing := something{value: "successful"}
	localCache.Set("new", testing)

	fmt.Println("set the variable")
	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	value, wasRetrieved := localCache.Get("new")
	if !wasRetrieved {
		fmt.Println("value not found")
	} else {
		fmt.Println("the found value: " + value.(something).value)
	}
}

func testWords() {
	cassConn := cassandra.Session{}
	cassConn.Init()
	recordUuid, _ := gocql.RandomUUID()
	// recordUuid, _ := gocql.ParseUUID("c92cf389-4451-4ccd-91d6-b20aed0fcf03")

	// accountUuid, _ := gocql.RandomUUID()
	accountUuid, _ := gocql.ParseUUID("4d2e9ace-474c-427f-a32d-cec835d1c688")

	recordToInsert := cassandra.WordsToRecords{Word: "done", AccountUuid: accountUuid, RecordUuid: []string{recordUuid.String()}}

	cassConn.UpsertWordsToRecords(recordToInsert)

	searchWords := []string{"adam", "testing"}
	found := cassConn.GetWordsToRecords(searchWords, accountUuid)

	for _, a := range found {
		fmt.Println(a.Word)
	}
	// SELECT * FROM words_to_records WHERE account_uuid = '4d2e9ace-474c-427f-a32d-cec835d1c688' and word IN ('adam','testing')
	cassConn.Close()
}