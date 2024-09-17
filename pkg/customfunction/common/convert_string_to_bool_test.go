package customFunc

import (
	"reflect"
	"testing"
)

func Test_ConvertString2Bool(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want FuncResult
	}{
		{"test ConvertString2Bool with str = YES",
			"YES",
			FuncResult{Result: true},
		},
		{"test ConvertString2Bool with str = Y",
			"YES",
			FuncResult{Result: true},
		},
		{"test ConvertString2Bool with str = yEs",
			"yEs",
			FuncResult{Result: true},
		},
		{"test ConvertString2Bool with str = no",
			"no",
			FuncResult{Result: false},
		},
		{"test ConvertString2Bool with str = empty",
			"",
			FuncResult{Result: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertString2Bool(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertString2Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}
