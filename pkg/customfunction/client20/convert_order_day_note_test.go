package funcClient20

import (
	"reflect"
	"testing"

	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
)

func Test_ConvertOrderDay(t *testing.T) {
	tests := []struct {
		name         string
		orderDayNote string
		want         customFunc.FuncResult
	}{
		{
			name:         "orderDayNote is empty, return empty OrderDay",
			orderDayNote: "",
			want: customFunc.FuncResult{
				Result: nil,
			},
		},
		{
			name:         "orderDayNote is EVERYDAY, return OrderDay with EveryDay is true",
			orderDayNote: "evErYdaY",
			want: customFunc.FuncResult{
				Result: OrderDay{
					EveryDay: true,
				},
			},
		},
		{
			name:         "orderDayNote is T2->6,CN, return OrderDay with Weekly is T2->6,8",
			orderDayNote: "T2->6,CN",
			want: customFunc.FuncResult{
				Result: OrderDay{
					Weekly: "2,3,4,5,6,8",
				},
			},
		},
		{
			name:         "orderDayNote is T2->6,SUN, return OrderDay with Weekly is T2->6,8",
			orderDayNote: "T2->6,SUN",
			want: customFunc.FuncResult{
				Result: OrderDay{
					Weekly: "2,3,4,5,6,8",
				},
			},
		},
		{
			name:         "orderDayNote is invalid position weekly, return error message",
			orderDayNote: "T2->6,3,4",
			want: customFunc.FuncResult{
				ErrorMessage: errorz.ErrFormatOrderScheduler.Error(),
			},
		},
		{
			name:         "orderDayNote is swap position weekly, return error message",
			orderDayNote: "T2->6,8->3",
			want: customFunc.FuncResult{
				ErrorMessage: errorz.ErrFormatOrderScheduler.Error(),
			},
		},
		{
			name:         "orderDayNote is T2,3,4->6, return OrderDay with Weekly is T2,3,4,5,6",
			orderDayNote: "T2,3,4->6",
			want: customFunc.FuncResult{
				Result: OrderDay{
					Weekly: "2,3,4,5,6",
				},
			},
		},
		{
			name:         "orderDayNote is duplicate weekly, return error",
			orderDayNote: "T2,3,4->6,2,3,4",
			want: customFunc.FuncResult{
				ErrorMessage: errorz.ErrFormatOrderScheduler.Error(),
			},
		},
		{
			name:         "orderDayNote is all week, return OrderDay with EveryDay is true",
			orderDayNote: "T2->8",
			want: customFunc.FuncResult{
				Result: OrderDay{
					EveryDay: true,
					Weekly:   "2,3,4,5,6,7,8",
				},
			},
		},
		{
			name:         "orderDayNote is all week with SUN, return OrderDay with EveryDay is true",
			orderDayNote: "T2->CN",
			want: customFunc.FuncResult{
				Result: OrderDay{
					EveryDay: true,
					Weekly:   "2,3,4,5,6,7,8",
				},
			},
		},
		{
			name:         "orderDayNote is N2->5,9->20,30, return OrderDay with Monthly is N2->5,9->20,30",
			orderDayNote: "N2->5,9->20,30",
			want: customFunc.FuncResult{
				Result: OrderDay{
					Monthly: "2,3,4,5,9,10,11,12,13,14,15,16,17,18,19,20,30",
				},
			},
		},
		{
			name:         "orderDayNote is duplicate monthly, return error",
			orderDayNote: "N2->5,10,10,11",
			want: customFunc.FuncResult{
				ErrorMessage: errorz.ErrFormatOrderScheduler.Error(),
			},
		},
		{
			name:         "orderDayNote is all month, return OrderDay with EveryDay is true",
			orderDayNote: "N1->31",
			want: customFunc.FuncResult{
				Result: OrderDay{
					EveryDay: true,
					Monthly:  "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31",
				},
			},
		},
		{
			name:         "orderDayNote is invalid, return error message",
			orderDayNote: "T2->6,CN,8",
			want: customFunc.FuncResult{
				ErrorMessage: errorz.ErrFormatOrderScheduler.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertOrderDay(tt.orderDayNote); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertOrderDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
