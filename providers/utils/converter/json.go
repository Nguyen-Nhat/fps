package converter

import (
	"reflect"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

// Override ... override common json by input json
func Override(input, common interface{}) {
	switch inputData := input.(type) {
	case []map[string]interface{}:
		switch commonData := common.(type) {
		case []interface{}:
			for idx, v := range inputData {
				Override(v, commonData[idx])
			}
		default:
			logger.Errorf("no support %v", commonData)
		}
	case map[string]interface{}:
		switch commonData := common.(type) {
		case map[string]interface{}:
			for k, v := range commonData {
				switch reflect.TypeOf(v).Kind() {
				case reflect.Slice:
					if inputData[k] == nil {
						inputData[k] = []interface{}{}
					}
					Override(inputData[k], v)
				case reflect.Map:
					if inputData[k] == nil {
						inputData[k] = map[string]interface{}{}
					}
					Override(inputData[k], v)
				default:
					// do simply replacement for primitive type
					_, ok := inputData[k]
					if !ok {
						inputData[k] = v
					}
				}
			}
		default:
			logger.Errorf("no support %v", commonData)
		}
	case []interface{}:
		switch commonData := common.(type) {
		case []interface{}:
			inputData = commonData
		}
	default:
		logger.Errorf("no support %v", inputData)
	}
}
