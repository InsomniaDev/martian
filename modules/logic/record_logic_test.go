package logic

import (
	"reflect"
	"testing"
)

func TestSortAndRetrieveRecordUuids(t *testing.T) {
	type args struct {
		records              [][]string
		numOfRecordsToReturn int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				records: [][]string{
					{
						"test", "response", "new",
					},
					{
						"test", "response",
					},
					{
						"response",
					},
				},
				numOfRecordsToReturn: 1,
			},
			want: []string{"response"},
		},
		{
			name: "test2",
			args: args{
				records: [][]string{
					{
						"test", "response", "new",
					},
					{
						"test", "response",
					},
					{
						"response",
					},
				},
				numOfRecordsToReturn: 2,
			},
			want: []string{"response","test"},
		},
		{
			name: "test3",
			args: args{
				records: [][]string{
					{
						"test", "response", "new",
					},
					{
						"test", "response",
					},
					{
						"response",
					},
					{
						"test", "response", "new",
					},
					{
						"test", "response",
					},
					{
						"response",
					},
					{
						"test", "response", "new",
					},
				},
				numOfRecordsToReturn: 2,
			},
			want: []string{"response","test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortAndRetrieveRecordUuids(tt.args.records, tt.args.numOfRecordsToReturn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortAndRetrieveRecordUuids() = %v, want %v", got, tt.want)
			}
		})
	}
}
