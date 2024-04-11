package funcClient20

import (
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const (
	everyDayUpperCase = "EVERYDAY"
	weeklyPrefixStr   = "T"
	monthlyPrefixStr  = "N"

	numberOfOrderSchedulerWeekly  = 7
	numberOfOrderSchedulerMonthly = 30

	startWeek  = 2
	endWeek    = 8
	startMonth = 1
	endMonth   = 31

	sundayNumber = "8"
	sundayVi     = "CN"
	sundayEn     = "SUN"
)

// ConvertOrderDay ...
func ConvertOrderDay(orderDayNote string) customFunc.FuncResult {
	orderDayNote = strings.ToUpper(orderDayNote)
	if orderDayNote == constant.EmptyString || orderDayNote == everyDayUpperCase {
		return customFunc.FuncResult{
			Result: OrderDay{
				EveryDay: true,
			},
		}
	}

	isEveryDay := false
	if strings.HasPrefix(orderDayNote, weeklyPrefixStr) {
		orderDayNote = strings.ReplaceAll(orderDayNote, sundayVi, sundayNumber)
		orderDayNote = strings.ReplaceAll(orderDayNote, sundayEn, sundayNumber)
		orderDays, err := convertOrderDayWeekMonth(orderDayNote, true)
		if err != nil {
			return customFunc.FuncResult{
				ErrorMessage: err.Error(),
			}
		}
		if len(orderDays) >= numberOfOrderSchedulerWeekly {
			isEveryDay = true
		}
		return customFunc.FuncResult{
			Result: OrderDay{
				EveryDay: isEveryDay,
				Weekly:   utils.JoinIntArray2String(orderDays, constant.SplitByComma),
			},
		}
	}
	if strings.HasPrefix(orderDayNote, monthlyPrefixStr) {
		orderDays, err := convertOrderDayWeekMonth(orderDayNote, false)
		if err != nil {
			return customFunc.FuncResult{
				ErrorMessage: err.Error(),
			}
		}
		if len(orderDays) > numberOfOrderSchedulerMonthly || isOrderDayAllMonth(orderDays) {
			isEveryDay = true
		}
		return customFunc.FuncResult{
			Result: OrderDay{
				EveryDay: isEveryDay,
				Monthly:  utils.JoinIntArray2String(orderDays, constant.SplitByComma),
			},
		}
	}
	return customFunc.FuncResult{
		ErrorMessage: errorz.ErrFormatOrderScheduler.Error(),
	}
}

func convertOrderDayWeekMonth(orderDayNote string, isWeekly bool) ([]int32, error) {
	orderDayNotes := strings.Split(orderDayNote[1:], constant.SplitByComma)
	lastOrderDay := int32(0)
	listOrderDay := make([]int32, 0)
	for _, orderDay := range orderDayNotes {
		orderDays, err := convertOrderDayDeduce(orderDay)
		if err != nil {
			return nil, err
		}
		if orderDays[0] <= lastOrderDay {
			return nil, errorz.ErrFormatOrderScheduler
		}
		lastOrderDay = orderDays[len(orderDays)-1]
		listOrderDay = append(listOrderDay, orderDays...)
	}
	if isWeekly {
		if listOrderDay[0] < startWeek {
			return nil, errorz.ErrFormatOrderScheduler
		}
		if listOrderDay[len(listOrderDay)-1] > endWeek {
			return nil, errorz.ErrFormatOrderScheduler
		}
		return listOrderDay, nil
	}
	if listOrderDay[0] < startMonth {
		return nil, errorz.ErrFormatOrderScheduler
	}
	if listOrderDay[len(listOrderDay)-1] > endMonth {
		return nil, errorz.ErrFormatOrderScheduler
	}
	return listOrderDay, nil
}

func isOrderDayAllMonth(orderDays []int32) bool {
	for idx, orderDay := range orderDays {
		if int(orderDay) != idx+1 {
			return false
		}
	}
	return true
}

func convertOrderDayDeduce(orderDay string) ([]int32, error) {
	orderDays, err := utils.String2ArrayInt32(orderDay, constant.SplitByDeduce)
	if err != nil {
		return nil, errorz.ErrFormatOrderScheduler
	}
	if len(orderDays) == constant.One {
		return orderDays, nil
	}
	if len(orderDays) != constant.Two {
		return nil, errorz.ErrFormatOrderScheduler
	}
	if orderDays[0] >= orderDays[1] {
		return nil, errorz.ErrFormatOrderScheduler
	}
	listOrderDay := make([]int32, orderDays[1]-orderDays[0]+1)
	for idx := orderDays[0]; idx <= orderDays[1]; idx++ {
		listOrderDay[idx-orderDays[0]] = idx
	}
	return listOrderDay, nil
}
