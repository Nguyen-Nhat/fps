package executerowgroup

import (
	"fmt"
	"reflect"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

func mergeMapInterface(first, second map[string]interface{}) (map[string]interface{}, error) {
	// 1. Init Result
	result := make(map[string]interface{})

	// 2. Copy First map
	for key, value := range first {
		result[key] = value
	}

	// 3. Merge Second map to Result
	for key, value := range second {
		existedValue, existed := result[key]

		// 3.1. If key hasn't appeared -> add to result
		if !existed {
			result[key] = value
			continue
		}

		// Merge value
		rfValue := reflect.ValueOf(value)
		rfExistedValue := reflect.ValueOf(existedValue)

		// Case Type Not Match
		if rfValue.Kind() != rfExistedValue.Kind() {
			logger.Errorf("mergeMapInterface() got error Type Not Match for key=%s", key)
			return nil, fmt.Errorf("cannot group value for field %s", key)
		}

		// Case is not Array
		if rfValue.Kind() != reflect.Slice && rfValue.Kind() != reflect.Array {
			if existedValue == value { // same value -> do nothing, go to next key-value
				continue
			} else { // value not same -> cannot merge
				logger.Errorf("mergeMapInterface() got error Value Not Match for key=%s", key)
				return nil, fmt.Errorf("cannot group value for field %s", key)
			}
		}

		// Case is Array
		if rfValue.Kind() == reflect.Slice || rfValue.Kind() == reflect.Array {
			result[key] = appendInterface(existedValue, value)
		}
	}

	return result, nil
}

func appendInterface(a, b interface{}) interface{} {
	return reflect.AppendSlice(reflect.ValueOf(a), reflect.ValueOf(b)).Interface()
}
