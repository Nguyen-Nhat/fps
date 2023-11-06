package customFunc

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	errMissingParameter  = "missing parameter"
	errFunctionNoSupport = "not support function"

	customFunctionPrefix = "$func."
	paramSeparator       = ";"
)

// ToCustomFunction ...
// Conditional: functionPattern has to pass IsCustomFunction(...) function
func ToCustomFunction(functionPattern string) (CustomFunction, error) {
	// 1. Check is function
	if !IsCustomFunction(functionPattern) {
		return CustomFunction{}, fmt.Errorf("%s is not function", functionPattern)
	}

	// 2. Split to get function name and parameter
	functionWithParams := strings.TrimPrefix(functionPattern, customFunctionPrefix)
	splitRes := strings.Split(functionWithParams, paramSeparator)
	return CustomFunction{functionPattern, splitRes[0], splitRes[1:], nil}, nil
}

// IsCustomFunction ...
// E.g: $func.sys.randomInt()
func IsCustomFunction(functionPattern string) bool {
	return strings.HasPrefix(functionPattern, customFunctionPrefix)
}

func ExecuteFunction(cf CustomFunction) (FuncResult, error) {
	switch cf.Name {
	case funcTestError:
		return testFuncError(), nil
	case funcTest:
		if len(cf.ParamsMapped) < 2 {
			return FuncResult{}, fmt.Errorf(errMissingParameter)
		} else {
			first, err1 := strconv.Atoi(cf.ParamsMapped[0])
			second, err2 := strconv.Atoi(cf.ParamsMapped[1])
			if err1 != nil || err2 != nil {
				return FuncResult{}, fmt.Errorf("%v or %v is not number", cf.ParamsMapped[0], cf.ParamsMapped[1])
			} else {
				return testFunc(first, second), nil
			}
		}
	case funcReUploadFile:
		if len(cf.ParamsMapped) < 1 {
			return FuncResult{}, fmt.Errorf(errMissingParameter)
		} else {
			return reUploadFile(cf.ParamsMapped[0]), nil
		}
	default:
		return FuncResult{}, fmt.Errorf(errFunctionNoSupport)
	}

}
