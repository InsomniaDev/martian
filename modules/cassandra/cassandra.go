package cassandra

import (
	"github.com/gocql/gocql"
)

type Session struct {
	Connection *gocql.Session
}

func (s *Session) Init() {
	var err error
	cluster := gocql.NewCluster("192.168.1.19:30506")
	cluster.Keyspace = "martian"
	s.Connection, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
}

func (s *Session) Close() {
	s.Connection.Close()
}

func (s *Session) ExecuteBatch(batch *gocql.Batch) (err error) {
	err = s.Connection.ExecuteBatch(batch)
	return
}
