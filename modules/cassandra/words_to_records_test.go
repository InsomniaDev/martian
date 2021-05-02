package cassandra

import (
	"testing"

	"github.com/gocql/gocql"
)

func TestSession_DeleteRecordsFromWords(t *testing.T) {
	type fields struct {
		Connection *gocql.Session
	}
	type args struct {
		words WordsToRecords
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Connection: tt.fields.Connection,
			}
			s.DeleteRecordsFromWords(tt.args.words)
		})
	}
}
