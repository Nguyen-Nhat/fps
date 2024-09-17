package customFunc

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// ConvertDateTimeFormat ...
// Convert date time format from current format to expect format
func ConvertDateTimeFormat(dateTimeStr, currentFormat, expectFormat string) FuncResult {
	dateTime, err := utils.StringToDateInUTC7Location(dateTimeStr, currentFormat)
	if err != nil {
		return FuncResult{ErrorMessage: errorz.ErrDateTimeFormat(dateTimeStr, currentFormat)}
	}
	return FuncResult{Result: dateTime.Format(expectFormat)}
}
