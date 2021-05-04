package logic

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/insomniadev/martian/modules/cassandra"
)

// TODO: Process an incoming request

func ParseInsert(conn *cassandra.Session, accountUuid gocql.UUID, ) {
	fmt.Println("I don't do anything yet")
	
}

func UpsertRecord(conn *cassandra.Session, record cassandra.Record) bool {
	return conn.UpsertRecord(record)
}