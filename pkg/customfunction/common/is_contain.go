package customFunc

import (
	"go.tekoapis.com/tekone/library/util/slice"
	"strings"
)

// IsContain ...
// Check str contains subString. Default is case-sensitive.
// Flag mustEqual require str and subStr are equal
func IsContain(str string, subStr string, isCaseInsensitive string, mustEqual string) FuncResult {
	if slice.Contains(listValueTrue, strings.ToUpper(mustEqual)) {
		return FuncResult{Result: strings.EqualFold(str, subStr)}
	}
	if slice.Contains(listValueTrue, strings.ToUpper(isCaseInsensitive)) {
		str = strings.ToLower(str)
		subStr = strings.ToLower(subStr)
	}
	return FuncResult{Result: strings.Contains(str, subStr)}
}
