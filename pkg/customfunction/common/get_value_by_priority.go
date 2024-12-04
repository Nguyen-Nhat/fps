package customFunc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/errorz"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

var (
	listValueTrue  = []string{"1", "TRUE", "T", "Y", "YES"}
	listValueFalse = []string{"0", "FALSE", "F", "N", "NO"}
)

// GetValueByPriority ...
// first params is type of response value
// if cant parse value to response type, return error
func GetValueByPriority(responseType string, values []string) FuncResult {
	valueStrMatched := constant.EmptyString
	isMatched := false
	for _, value := range values {
		if value == constant.EmptyString || value == constant.NilString {
			continue
		}
		valueStrMatched = value
		isMatched = true
		break
	}
	if !isMatched {
		return FuncResult{Result: nil}
	}

	switch responseType {
	case constant.TypeString:
		return FuncResult{Result: valueStrMatched}
	case constant.TypeInteger:
		if valueInt64, err := strconv.ParseInt(valueStrMatched, 10, 64); err == nil {
			return FuncResult{Result: valueInt64}
		} else {
			return FuncResult{Result: nil, ErrorMessage: errorz.ErrCantParseValue(valueStrMatched, constant.TypeInteger)}
		}
	case constant.TypeNumber:
		if valueFloat64, err := strconv.ParseFloat(valueStrMatched, 64); err == nil {
			return FuncResult{Result: valueFloat64}
		} else {
			return FuncResult{Result: nil, ErrorMessage: errorz.ErrCantParseValue(valueStrMatched, constant.TypeNumber)}
		}
	case constant.TypeBoolean:
		if valueBool, err := parseBool(valueStrMatched); err == nil {
			return FuncResult{Result: valueBool}
		} else {
			return FuncResult{Result: nil, ErrorMessage: errorz.ErrCantParseValue(valueStrMatched, constant.TypeBoolean)}
		}
	case constant.TypeJson:
		result := gjson.Parse(valueStrMatched)
		return FuncResult{Result: result.Value()}
	}

	return FuncResult{Result: nil}
}

func parseBool(valueStr string) (bool, error) {
	valueStr = strings.TrimSpace(valueStr)
	if utils.Contains(listValueTrue, strings.ToUpper(valueStr)) {
		return true, nil
	}
	if utils.Contains(listValueFalse, strings.ToUpper(valueStr)) {
		return false, nil
	}
	return false, fmt.Errorf("invalid value for bool: %s", valueStr)
}
