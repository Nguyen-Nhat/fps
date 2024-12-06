package customFunc

import (
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// IsContain ...
// Check str contains subString. Default is case-sensitive.
// Flag mustEqual require str and subStr are equal
func IsContain(str string, subStr string, isCaseInsensitive string, mustEqual string) FuncResult {
	if utils.Contains(listValueTrue, strings.ToUpper(mustEqual)) {
		return FuncResult{Result: strings.EqualFold(str, subStr)}
	}
	if utils.Contains(listValueTrue, strings.ToUpper(isCaseInsensitive)) {
		str = strings.ToLower(str)
		subStr = strings.ToLower(subStr)
	}
	return FuncResult{Result: strings.Contains(str, subStr)}
}
