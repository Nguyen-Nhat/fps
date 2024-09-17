package customFunc

import (
	"reflect"
	"testing"
)

func Test_ConvertDateTimeFormat(t *testing.T) {
	tests := []struct {
		name          string
		dateTimeStr   string
		currentFormat string
		expectFormat  string
		want          FuncResult
	}{
		{"test ConvertDateTimeFormat from split by slash to split by dash",
			"30/06/2024",
			"02/01/2006",
			"02-01-2006",
			FuncResult{Result: "30-06-2024"},
		},
		{"test ConvertDateTimeFormat from split by dash to split by slash",
			"30-06-2024",
			"02-01-2006",
			"02/01/2006",
			FuncResult{Result: "30/06/2024"},
		},
		{"test ConvertDateTimeFormat from split by slash to split by dash with time",
			"30/06/2024",
			"02/01/2006",
			"02-01-2006 15:04:05",
			FuncResult{Result: "30-06-2024 00:00:00"},
		},
		{"test ConvertDateTimeFormat from split by dash to split by slash with time",
			"30-06-2024",
			"02-01-2006",
			"02/01/2006 15:04:05",
			FuncResult{Result: "30/06/2024 00:00:00"},
		},
		{"test ConvertDateTimeFormat from split by dash to split by slash with time and other format",
			"30-06-2024",
			"02-01-2006",
			"2006-01-02T15:04:05Z",
			FuncResult{Result: "2024-06-30T00:00:00Z"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertDateTimeFormat(tt.dateTimeStr, tt.currentFormat, tt.expectFormat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertString2Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}
