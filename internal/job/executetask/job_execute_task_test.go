package executetask

import (
	"reflect"
	"testing"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessingrow"
)

func Test_toResponseResult(t *testing.T) {
	type args struct {
		task               fileprocessingrow.ProcessingFileRow
		taskMappingUpdated string
		curl               string
		responseBody       string
		messageRes         string
		status             int16
		startAt            time.Time
	}
	startedAt := time.Now()
	tests := []struct {
		name string
		args args
		want fileprocessingrow.UpdateAfterExecutingByJob
	}{
		{"",
			args{
				task:               fileprocessingrow.ProcessingFileRow{},
				taskMappingUpdated: "taskMappingUpdated",
				curl:               "curl",
				responseBody:       "responseBody",
				messageRes:         "messageRes",
				status:             1,
				startAt:            startedAt,
			},
			fileprocessingrow.UpdateAfterExecutingByJob{
				Task:         fileprocessingrow.ProcessingFileRow{},
				TaskMapping:  "taskMappingUpdated",
				RequestCurl:  "curl",
				ResponseRaw:  "responseBody",
				Status:       1,
				ErrorDisplay: "messageRes",
				ExecutedTime: 0,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toResponseResult(tt.args.task, tt.args.taskMappingUpdated, tt.args.curl, tt.args.responseBody, tt.args.messageRes, tt.args.status, tt.args.startAt)
			got.ExecutedTime = 0

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toResponseResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
