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

func Test_IsColumnIndex(t *testing.T) {
	tests := []struct {
		name       string
		columnName string
		want       bool
	}{
		// case invalid -> empty
		{"test IsColumnIndex: when columnName none then return false",
			"", false},
		{"test IsColumnIndex: when columnName has 1 character, missing prefix then return false",
			"A", false},
		{"test IsColumnIndex: when columnName has 3 character, missing prefix then return false",
			"ABC", false},
		{"test IsColumnIndex: when columnName has only prefix then return false",
			"$", false},
		{"test IsColumnIndex: when columnName has 2 character, missing prefix then return false",
			"A$", false},
		{"test IsColumnIndex: when columnName has 3 character but one of them is not from A-Z then return false",
			"$A1", false},
		// case valid -> return correct value
		{"test IsColumnIndex: when columnName has 2 character, start with prefix then return true",
			"$A", true},
		{"test IsColumnIndex: when columnName has 3 character, start with prefix then return true",
			"$AB", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsColumnIndex(tt.columnName); got != tt.want {
				t.Errorf("IsColumnIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
