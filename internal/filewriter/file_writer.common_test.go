package filewriter

import (
	"testing"
)

func TestWriteDataToCsv(t *testing.T) {
	type args struct {
		fileData [][]string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test WriteDataToCsv: case no data",
			args{nil},
		},
		{"test WriteDataToCsv: happy case",
			args{[][]string{
				{"a", "b", "c"},
				{"1", "2", "3"},
			},
			},
		},
		{"test WriteDataToCsv: happy case with empty string",
			args{[][]string{
				{"a", "b", "c", ""},
				{"1", "", "3"},
				{"1", "", "3", "", ""},
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := writeDataToCsv(tt.args.fileData)
			if err != nil {
				t.Errorf("writeDataToCsv() error = %v", err)
				return
			}
		})
	}
}

func TestWriteDataToExcel(t *testing.T) {
	type args struct {
		fileData  [][]string
		sheetName string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test WriteDataToExcel: case no data",
			args{nil, "sheetName"},
		},
		{"test WriteDataToExcel: happy case",
			args{[][]string{
				{"a", "b", "c"},
				{"1", "2", "3"},
			}, "sheetName",
			},
		},
		{"test WriteDataToExcel: happy case with empty string",
			args{[][]string{
				{"a", "b", "c", ""},
				{"1", "", "3"},
				{"1", "", "3", "", ""},
			}, "sheetName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := writeDataToXlsx(tt.args.fileData, tt.args.sheetName)
			if err != nil {
				t.Errorf("writeDataToXlsx() error = %v", err)
				return
			}
		})
	}
}
