package customFunc

import (
	"encoding/json"
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

// GetValueByPriorityAndMapping ...
// first params is dictionary. If it has value, return value of key in dictionary first
// second params is response type. Parse value with type
// third params is list of values
func GetValueByPriorityAndMapping(dict, responseType string, strs []string) FuncResult {
	valueStrMatched := constant.EmptyString
	isMatched := false
	for _, value := range strs {
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

	if dict != constant.EmptyString {
		if value, err := getValueDictionary(valueStrMatched, dict); err == nil {
			valueStrMatched = value
		} else {
			return FuncResult{ErrorMessage: err.Error()}
		}
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

func getValueDictionary(key, dictionary string) (string, error) {
	var dict map[string]string
	err := json.Unmarshal([]byte(dictionary), &dict)
	if err != nil {
		return constant.EmptyString, err
	}
	keys := make([]string, 0)

	for k, v := range dict {
		if strings.EqualFold(k, key) {
			return v, nil
		}
		keys = append(keys, k)
	}
	return constant.EmptyString, errorz.ErrNotExistValueInList(key, keys)
}
