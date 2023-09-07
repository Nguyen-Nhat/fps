package converter

import (
	"reflect"
	"testing"
)

func Test_IntArrToInt32Arr(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int32
	}{
		{"test IntArrToInt32Arr case arr is nil -> empty array",
			nil, []int32{}},
		{"test IntArrToInt32Arr case arr is empty -> empty array",
			[]int{}, []int32{}},
		{"test IntArrToInt32Arr case arr has data -> correct array",
			[]int{-1, 0, 1, 11, 999999999, -999999999},
			[]int32{-1, 0, 1, 11, 999999999, -999999999}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntArrToInt32Arr(tt.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntArrToInt32Arr() = %v, want %v", got, tt.want)
			}
		})
	}
}
