package filewriter

import (
	"encoding/csv"
	"reflect"
	"testing"

	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
)

func Test_excelFileWriter_GetFileBytes(t *testing.T) {
	tests := []struct {
		name           string
		fileData       [][]string
		outputFileType string
		wantErr        bool
	}{
		{"test GetFileBytes: case no data", [][]string{}, constant.ExtFileXLSX, false},
		{"test GetFileBytes: happy case output csv", [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}, constant.ExtFileCSV, false},
		{"test GetFileBytes: happy case output xlsx", [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}, constant.ExtFileXLSX, false},
		{"test GetFileBytes: happy case output default", [][]string{
			{"a", "b", "c"},
			{"1", "2", "3"},
		}, "abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &excelFileWriter{sheetName: "data", outputFileType: tt.outputFileType}
			c.exFile = excelize.NewFile()
			_, _ = c.exFile.NewSheet(c.sheetName)
			for rowIndex, rowData := range tt.fileData {
				cell, _ := excelize.CoordinatesToCellName(1, rowIndex+1)
				_ = c.exFile.SetSheetRow(c.sheetName, cell, &rowData)
			}

			var gotData [][]string
			var err error
			if tt.outputFileType == constant.ExtFileCSV {
				got, err := c.GetFileBytes()
				if (err != nil) != tt.wantErr {
					t.Errorf("GetFileBytes() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				gotData, err = csv.NewReader(got).ReadAll()
			} else {
				gotData, err = c.exFile.GetRows(c.sheetName)
			}

			if err != nil || !reflect.DeepEqual(gotData, tt.fileData) {
				t.Errorf("GetFileBytes() got = %v, want %v", gotData, tt.fileData)
			}
		})
	}
}
