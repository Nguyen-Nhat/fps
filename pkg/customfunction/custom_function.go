package customFunc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	funcClient10 "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/client10"
	funcClient9 "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/client9"
	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

const (
	errMissingParameter  = "missing parameter"
	errFunctionNoSupport = "not support function"

	customFunctionPrefix = "$func."
	paramSeparator       = ";"
)

// ToCustomFunction ...
// Conditional: functionPattern has to pass IsCustomFunction(...) function
func ToCustomFunction(functionPattern string) (customFunc.CustomFunction, error) {
	// 1. Check is function
	if !IsCustomFunction(functionPattern) {
		return customFunc.CustomFunction{}, fmt.Errorf("%s is not function", functionPattern)
	}

	// 2. Split to get function name and parameter
	functionWithParams := strings.TrimPrefix(functionPattern, customFunctionPrefix)
	splitRes := strings.Split(functionWithParams, paramSeparator)
	return customFunc.CustomFunction{FunctionPattern: functionPattern, Name: splitRes[0], ParamsRaw: splitRes[1:]}, nil
}

// IsCustomFunction ...
// E.g: $func.sys.randomInt()
func IsCustomFunction(functionPattern string) bool {
	return strings.HasPrefix(functionPattern, customFunctionPrefix)
}

func ExecuteFunction(cf customFunc.CustomFunction) (customFunc.FuncResult, error) {
	startTime := time.Now()
	result, err := executeFunction(cf)
	logger.InfoT("ExecuteFunction > `%s`, params=%v -> %d (ms)", cf.Name, cf.ParamsMapped, time.Since(startTime))
	return result, err
}

func executeFunction(cf customFunc.CustomFunction) (customFunc.FuncResult, error) {
	switch cf.Name {
	case customFunc.FuncTestError:
		return customFunc.TestFuncError(), nil
	case customFunc.FuncTest:
		if len(cf.ParamsMapped) < 2 {
			return customFunc.FuncResult{}, fmt.Errorf(errMissingParameter)
		} else {
			first, err1 := strconv.Atoi(cf.ParamsMapped[0])
			second, err2 := strconv.Atoi(cf.ParamsMapped[1])
			if err1 != nil || err2 != nil {
				return customFunc.FuncResult{}, fmt.Errorf("%v or %v is not number", cf.ParamsMapped[0], cf.ParamsMapped[1])
			} else {
				return customFunc.TestFunc(first, second), nil
			}
		}
	case funcClient9.FuncReUploadFile:
		if len(cf.ParamsMapped) < 1 {
			return customFunc.FuncResult{}, fmt.Errorf(errMissingParameter)
		} else {
			return funcClient9.ReUploadFile(cf.ParamsMapped[0]), nil
		}
	case funcClient10.FuncConvertSellerSkuAndUomName:
		if len(cf.ParamsMapped) < 1 {
			return customFunc.FuncResult{}, fmt.Errorf(errMissingParameter)
		} else {
			return funcClient10.ConvertSellerSkus(cf.ParamsMapped[0]), nil
		}
	default:
		return customFunc.FuncResult{}, fmt.Errorf(errFunctionNoSupport)
	}

}
