package customFunc

import (
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const (
	yesStr = "YES"
	yStr   = "Y"
)

var (
	yesList = []string{yesStr, yStr}
)

// ConvertString2Bool ...
// Convert string to boolean, return true if string is "YES" or "Y", otherwise return false
func ConvertString2Bool(str string) FuncResult {
	if utils.Contains(yesList, strings.ToUpper(str)) {
		return FuncResult{Result: true}
	}
	return FuncResult{Result: false}
}
