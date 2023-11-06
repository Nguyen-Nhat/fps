package excel

import "testing"

func Test_getValueFromColumnKey(t *testing.T) {
	type args struct {
		columnKey string
		data      []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// case invalid -> empty
		{"test getValueFromColumnKey: when columnKey none then return empty",
			args{"", []string{"value_A", "value_B"}}, ""},
		{"test getValueFromColumnKey: when columnKey > 'Z' then return empty",
			args{"AA", []string{"value_A", "value_B"}}, ""},
		{"test getValueFromColumnKey: when columnKey wrong then return empty",
			args{".A", []string{"value_A", "value_B"}}, ""},
		{"test getValueFromColumnKey: when columnKey wrong then return empty",
			args{"A.", []string{"value_A", "value_B"}}, ""},
		{"test getValueFromColumnKey: when columnKey wrong then return empty",
			args{"A.asdlkfj", []string{"value_A", "value_B"}}, ""},
		{"test getValueFromColumnKey: when columnKey not existed then return empty",
			args{"Z", []string{"value_A", "value_B"}}, ""},
		// case valid -> return correct value
		{"test getValueFromColumnKey: when columnKey is A then return correct value",
			args{"A", []string{"value_A", "value_B"}}, "value_A"},
		{"test getValueFromColumnKey: when columnKey is B then return correct value",
			args{"B", []string{"value_A", "value_B"}}, "value_B"},
		{"test getValueFromColumnKey: when columnKey is C then return correct value",
			args{"C", []string{"value_A", "value_B", "", "value_D"}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetValueFromColumnKey(tt.args.columnKey, tt.args.data); got != tt.want {
				t.Errorf("getValueFromColumnKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
