package basejobmanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
)

const (
	errTypeWrong      = "sai kiểu dữ liệu"
	errTypeNotSupport = "không hỗ trợ kiểu dữ liệu"
)

// ConvertToRealValue ...
// note: return nil when valueStr is empty
func ConvertToRealValue(fieldType string, valueStr string, dependsOnKey string) (interface{}, error) {
	if len(valueStr) == 0 {
		return nil, nil
	}

	var realValue interface{}
	switch strings.ToLower(fieldType) {
	case configloader.TypeString:
		realValue = valueStr
	case configloader.TypeInteger:
		if valueInt64, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			realValue = valueInt64
		} else {
			return nil, fmt.Errorf("%s (%s)", errTypeWrong, dependsOnKey)
		}
	case configloader.TypeNumber:
		if valueFloat64, err := strconv.ParseFloat(valueStr, 64); err == nil {
			realValue = valueFloat64
		} else {
			return nil, fmt.Errorf("%s (%s)", errTypeWrong, dependsOnKey)
		}
	case configloader.TypeBoolean:
		if valueBool, err := strconv.ParseBool(valueStr); err == nil {
			realValue = valueBool
		} else {
			return nil, fmt.Errorf("%s (%s)", errTypeWrong, dependsOnKey)
		}
	case configloader.TypeJson:
		result := gjson.Parse(valueStr)
		realValue = result.Value()
	default:
		return nil, fmt.Errorf("%s %s", errTypeNotSupport, fieldType)
	}
	return realValue, nil
}
