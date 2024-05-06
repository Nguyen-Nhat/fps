package filewriter

import (
	"encoding/csv"
	"reflect"
	"testing"

	"github.com/xuri/excelize/v2"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
)

func Test_xlsFileWriter_GetFileBytes(t *testing.T) {
	tests := []struct {
		name           string
		fileData       [][]string
		outputFileType string
		wantErr        bool
	}{
		{"test GetFileBytes: case no data", nil, constant.ExtFileCSV, false},
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
			c := &xlsFileWriter{fileData: tt.fileData, outputFileType: tt.outputFileType}

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
				got, err := c.GetFileBytes()
				if (err != nil) != tt.wantErr {
					t.Errorf("GetFileBytes() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				xlsxFile, err := excelize.OpenReader(got)
				if err != nil {
					t.Errorf("GetFileBytes() error = %v", err)
					return
				}
				gotData, err = xlsxFile.GetRows(xlsxFile.GetSheetName(xlsxFile.GetActiveSheetIndex()))
			}

			if err != nil || !reflect.DeepEqual(gotData, tt.fileData) {
				t.Errorf("GetFileBytes() got = %v, want %v", gotData, tt.fileData)
			}
		})
	}
}

func Test_xlsFileWriter_UpdateDataInColumnOfFile(t *testing.T) {
	type fields struct {
		fileData       [][]string
		dataIndexStart int
	}
	type args struct {
		columnName string
		columnData map[int]string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantFileData [][]string
		wantErr      bool
	}{
		{"test UpdateDataInColumnOfFile: case columnName is columnIndex but wrong -> update first column",
			fields{[][]string{{"header_A", "header_B"}, {"value_A2", "value_B2"}}, 2}, args{"$AAAAAA", map[int]string{0: "value_A2_updated"}},
			[][]string{{"header_A", "header_B"}, {"value_A2_updated", "value_B2"}}, false},

		{"test UpdateDataInColumnOfFile: case columnName is columnIndex but wrong -> update correctly",
			fields{[][]string{{"header_A", "header_B"}, {"value_A2", "value_B2"}}, 2}, args{"$B", map[int]string{0: "value_B2_updated"}},
			[][]string{{"header_A", "header_B"}, {"value_A2", "value_B2_updated"}}, false},

		{"test UpdateDataInColumnOfFile: case columnName is columnIndex but wrong -> update correctly",
			fields{[][]string{{"header_A", "header_B"}, {"value_A2", "value_B2"}}, 2}, args{"$E", map[int]string{0: "value_E2"}},
			[][]string{{"header_A", "header_B", "", "", ""}, {"value_A2", "value_B2", "", "", "value_E2"}}, false},

		{"test UpdateDataInColumnOfFile: case columnName is headerName but not found -> update first column",
			fields{[][]string{{"header_A", "header_B"}, {"value_A2", "value_B2"}}, 2}, args{"header_Z", map[int]string{0: "value_A2_updated"}},
			[][]string{{"header_A", "header_B"}, {"value_A2_updated", "value_B2"}}, false},

		{"test UpdateDataInColumnOfFile: case columnName is existed headerName -> update correctly",
			fields{[][]string{{"header_A", "header_B"}, {"value_A2", "value_B2"}}, 2}, args{"header_B", map[int]string{0: "value_B2_updated"}},
			[][]string{{"header_A", "header_B"}, {"value_A2", "value_B2_updated"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &xlsFileWriter{
				fileData:       tt.fields.fileData,
				dataIndexStart: tt.fields.dataIndexStart,
			}
			if err := c.UpdateDataInColumnOfFile(tt.args.columnName, tt.args.columnData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDataInColumnOfFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			fileDataAfterUpdate := c.fileData
			if !reflect.DeepEqual(fileDataAfterUpdate, tt.wantFileData) {
				t.Errorf("GetFileBytes() got = %v, want %v", fileDataAfterUpdate, tt.wantFileData)
			}
		})
	}
}
