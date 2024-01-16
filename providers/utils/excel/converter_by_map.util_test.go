package excel

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
)

func Test_mappingHeader(t *testing.T) {
	tests := []struct {
		name        string
		columnsData []dto.CellData[string]
		want        map[string]int
		wantErr     bool
	}{
		{"test happy", []dto.CellData[string]{
			{ColumnName: "A"},
			{ColumnName: "B"},
			{ColumnName: "Z"},
			{ColumnName: "AA"},
			{ColumnName: "AZ"},
		}, map[string]int{"A": 0, "B": 1, "Z": 25, "AA": 26, "AZ": 51}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mappingHeader(tt.columnsData)
			if (err != nil) != tt.wantErr {
				t.Errorf("mappingHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mappingHeaders() got = %v, want %v", got, tt.want)
			}
		})
	}
}
