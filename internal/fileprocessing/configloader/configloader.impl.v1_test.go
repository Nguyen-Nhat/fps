package configloader

import (
	"reflect"
	"testing"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/configtask"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
)

func Test_toRowGroupMD(t *testing.T) {
	tests := []struct {
		name string
		task configtask.ConfigTask
		want RowGroupMD
	}{
		{"test toRowGroupMD when config task NOT have group config -> return data has empty GroupByColumns", mockConfigTask("", 0),
			RowGroupMD{"", []int{}, 0}},
		{"test toRowGroupMD when config task have group config -> data with correct all fields", mockConfigTask("A,B,C", 100),
			RowGroupMD{"A,B,C", []int{0, 1, 2}, 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toRowGroupMD(tt.task); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toRowGroupMD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockConfigTask(groupByColumns string, groupBySizeLimit int32) configtask.ConfigTask {
	return configtask.ConfigTask{
		ConfigTask: ent.ConfigTask{
			GroupByColumns:   groupByColumns,
			GroupBySizeLimit: groupBySizeLimit,
		},
	}
}
