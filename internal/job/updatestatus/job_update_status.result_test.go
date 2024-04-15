package updatestatus

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
)

func Test_toTaskIDs(t *testing.T) {
	resultFileConfigs := []configloader.ResultFileConfigMD{
		{ValueInTaskID: 1},
		{ValueInTaskID: 1},
		{ValueInTaskID: 2},
		{ValueInTaskID: 3},
	}
	expectTaskIDs := []int32{1, 2, 3}

	if got := toTaskIDs(resultFileConfigs); !reflect.DeepEqual(got, expectTaskIDs) {
		t.Errorf("toTaskIDs() = %v, want %v", got, expectTaskIDs)
	}
}

func Test_getColumnData(t *testing.T) {
	resCfgTask1 := configloader.ResultFileConfigMD{ValuePath: "processingResult", ValueInTaskID: 1}
	resCfgTask2 := configloader.ResultFileConfigMD{ValuePath: "errorMessage", ValueInTaskID: 2}

	mockResultAsyncDAO := []fileprocessingrow.ResultAsyncDAO{
		{0, 1, `{"processingResult": "SUCCESS_0_1", "errorMessage": "NONE_0_1"}`},
		{0, 2, `{"processingResult": "FAILED_0_2", "errorMessage": "INVALID_DATA_0_2"}`},
		{1, 1, `{"processingResult": "SUCCESS_1_1", "errorMessage": "NONE_1_1"}`},
		{1, 2, `{"processingResult": "FAILED_1_2", "errorMessage": "INVALID_DATA_1_2"}`},
	}

	type args struct {
		resultFileConfig configloader.ResultFileConfigMD
		results          []fileprocessingrow.ResultAsyncDAO
	}
	tests := []struct {
		name string
		args args
		want map[int]string
	}{
		{"test getColumnData: case valuePath=processingResult and task=1", args{resCfgTask1, mockResultAsyncDAO},
			map[int]string{0: "SUCCESS_0_1", 1: "SUCCESS_1_1"}},
		{"test getColumnData: case valuePath=errorMessage and task=2", args{resCfgTask2, mockResultAsyncDAO},
			map[int]string{0: "INVALID_DATA_0_2", 1: "INVALID_DATA_1_2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getColumnData(tt.args.results, tt.args.resultFileConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getColumnData() = %v, want %v", got, tt.want)
			}
		})
	}
}
