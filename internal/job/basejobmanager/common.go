package basejobmanager

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/fileprocessing/configloader"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/tools/i18n"
)

// ConvertToRealValue ...
// note: return nil when valueStr is empty
func ConvertToRealValue(ctx context.Context, fieldType string, valueStr string, dependsOnKey string) (interface{}, error) {
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
			return nil, fmt.Errorf(i18n.GetMessageCtx(ctx, "errTypeWrong", "name", dependsOnKey))
		}
	case configloader.TypeNumber:
		if valueFloat64, err := strconv.ParseFloat(valueStr, 64); err == nil {
			realValue = valueFloat64
		} else {
			logger.Debugf("wrong data type: expect type %s but data is %+v", fieldType, valueStr)
			return nil, fmt.Errorf(i18n.GetMessageCtx(ctx, "errTypeWrong", "name", dependsOnKey))
		}
	case configloader.TypeBoolean:
		if valueBool, err := strconv.ParseBool(valueStr); err == nil {
			realValue = valueBool
		} else {
			return nil, fmt.Errorf(i18n.GetMessageCtx(ctx, "errTypeWrong", "name", dependsOnKey))
		}
	case configloader.TypeJson:
		result := gjson.Parse(valueStr)
		realValue = result.Value()
	default:
		return nil, fmt.Errorf(i18n.GetMessageCtx(ctx, "errTypeNotSupport", "name", fieldType))
	}
	return realValue, nil
}
