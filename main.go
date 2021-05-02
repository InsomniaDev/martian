package main

import (
	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

func main() {
	cassConn := cassandra.Session{}
	cassConn.Init()
	// recordUuid, _ := gocql.RandomUUID()
	recordUuid, _ := gocql.ParseUUID("c92cf389-4451-4ccd-91d6-b20aed0fcf03")
	
	// accountUuid, _ := gocql.RandomUUID()
	accountUuid, _ := gocql.ParseUUID("4d2e9ace-474c-427f-a32d-cec835d1c688")

	recordToInsert := cassandra.WordsToRecords{Word: "testing", AccountUuid: accountUuid, RecordUuid: []string{recordUuid.String()}}

	cassConn.DeleteRecordsFromWords(recordToInsert)
	cassConn.Close()
}
