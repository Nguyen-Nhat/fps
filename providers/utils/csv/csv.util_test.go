package csv

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readCsv(t *testing.T) {
	tests := []struct {
		name    string
		pathCsv string
		want    [][]string
	}{
		{
			name:    "test get csv happy case",
			pathCsv: "./test_data/happycase.csv",
			want: [][]string{
				{"header1,", "header2,"},
				{"a,b,c,", "a-b-c-"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.pathCsv)
			assert.Nil(t, err)
			defer file.Close()
			got, err := csv.NewReader(file).ReadAll()
			assert.Nil(t, err)
			for i, arr := range tt.want {
				for j, str := range arr {
					assert.Equal(t, str, got[i][j])
				}
			}
		})
	}
}
