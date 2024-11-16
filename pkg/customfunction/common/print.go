package customFunc

import (
	"fmt"
)

// PrintFormat ...
// return string format by fmt.Sprintf
func PrintFormat(format string, params []string) FuncResult {
	anyParams := make([]any, len(params))
	for i, v := range params {
		anyParams[i] = v
	}
	return FuncResult{Result: fmt.Sprintf(format, anyParams...)}
}
