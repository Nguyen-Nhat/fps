package filewriter

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/configmapping"
)

func Test_NewFileWriter(t *testing.T) {
	type args struct {
		fileURL           string
		sheetName         string
		dataIndexStart    int
		inputFileType     string
		outputFileTypeCfg configmapping.OutputFileType
	}
	tests := []struct {
		name       string
		args       args
		wantStruct *csvFileWriter
		wantErr    bool
	}{
		{"test NewFileWriter: case output not support -> error",
			args{"https://a.com", "sheetName", 2, constant.ExtFileXLSX, configmapping.OutputFileType("abc")}, nil, true},

		{"test NewFileWriter: case CSV, wrong URL -> error",
			args{"https://a.com", "sheetName", 2, constant.ExtFileCSV, configmapping.OutputFileTypeCSV}, nil, true},
		{"test NewFileWriter: case XLSX, wrong URL -> error",
			args{"https://a.com", "sheetName", 2, constant.ExtFileXLSX, configmapping.OutputFileTypeXLSX}, nil, true},
		{"test NewFileWriter: case XLS, wrong URL -> error",
			args{"https://a.com", "sheetName", 2, constant.ExtFileXLS, configmapping.OutputFileTypeXLS}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileWriter(tt.args.fileURL, tt.args.sheetName, tt.args.dataIndexStart, tt.args.inputFileType, tt.args.outputFileTypeCfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(got) != reflect.TypeOf(tt.wantStruct) {
				t.Errorf("NewFileWriter() got = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.wantStruct))
			}
		})
	}

}
