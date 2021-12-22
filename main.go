package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/brain"
	"github.com/insomniadev/martian/graphql"
	"github.com/insomniadev/martian/modules/cache"
	"github.com/insomniadev/martian/modules/cassandra"
)

var mainBrain brain.Brain

func subscriptionSubscriber(channel, payload string) {

	payload, _ = strconv.Unquote(payload)
	eventData := brain.Event{}
	_ = json.Unmarshal([]byte(payload), &eventData)
	eventData.Time = time.Now()
	mainBrain.ProcessEvent(eventData)
}

func main() {
	// testLocalCache()
	graphql.Graphql()

	// mainBrain.Init()
	// redispub.NewSubscriber("brain", subscriptionSubscriber)
	// harmony.RetrieveAllNodes()
	// go graphql.Graphql()
	// server.Start()
}

func testLocalCache() {
	localCache := cache.LocalCache{}
	localCache.Init()

	type something struct {
		value string
	}
	testing := something{value: "successful"}
	localCache.Set("new", testing)

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	value, wasRetrieved := localCache.Get("new")
	if !wasRetrieved {
		log.Println("value not found")
	} else {
		log.Println("the found value: " + value.(something).value)
	}
}

func testWords() {
	cassConn := cassandra.Session{}
	cassConn.Init()
	recordUuid, _ := gocql.RandomUUID()
	// recordUuid, _ := gocql.ParseUUID("c92cf389-4451-4ccd-91d6-b20aed0fcf03")

	// accountUuid, _ := gocql.RandomUUID()
	accountUuid, _ := gocql.ParseUUID("4d2e9ace-474c-427f-a32d-cec835d1c688")

	recordToInsert := cassandra.WordsToRecord{Word: "done", AccountUuid: accountUuid, RecordUuid: recordUuid.String()}

	cassConn.UpsertWordsToRecord(recordToInsert)

	searchWords := []string{"adam", "testing"}
	found := cassConn.GetWordsToRecords(searchWords, accountUuid)

	for _, a := range found {
		log.Println(a.Word)
	}
	// SELECT * FROM words_to_records WHERE account_uuid = '4d2e9ace-474c-427f-a32d-cec835d1c688' and word IN ('adam','testing')
	cassConn.Close()
}

// Old Main File for Rusty

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"time"

// 	"github.com/insomniadev/rusty/brain"
// 	graphql "github.com/insomniadev/rusty/graphql"
// 	"github.com/insomniadev/rusty/integrations/harmony"
// 	"github.com/insomniadev/rusty/modules/redispub"
// )

// var mainBrain brain.Brain

// func subscriptionSubscriber(channel, payload string) {

// 	payload, err := strconv.Unquote(payload)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	eventData := brain.Event{}
// 	err = json.Unmarshal([]byte(payload), &eventData)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	eventData.Time = time.Now()
// 	mainBrain.ProcessEvent(eventData)
// }

// func main() {
// 	mainBrain.Init()
// 	redispub.NewSubscriber("brain", subscriptionSubscriber)
// 	harmony.RetrieveAllNodes()
// 	graphql.Graphql()
// }
