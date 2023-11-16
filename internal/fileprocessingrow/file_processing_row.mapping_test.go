package fileprocessingrow

import (
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

func Test_getTaskName(t *testing.T) {
	taskMapping := `{
						"Tasks":[
							{"TaskIndex":1,"TaskName":"this is task 1"},
							{"TaskIndex":2,"TaskName":"this is task 2"},
							{"TaskIndex":3,"TaskName":"this is task 3"}
						]
					}`
	tests := []struct {
		name      string
		taskIndex int32
		want      string
	}{
		{"test getTaskName case task 1", 1, "this is task 1"},
		{"test getTaskName case task 2", 2, "this is task 2"},
		{"test getTaskName case task 3", 3, "this is task 3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfr := &ProcessingFileRow{ent.ProcessingFileRow{TaskIndex: tt.taskIndex, TaskMapping: taskMapping}}
			if got := getTaskName(pfr); got != tt.want {
				t.Errorf("getTaskName() = %v, want %v", got, tt.want)
			}
		})
	}
}
