package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

// BindRequestParamsToStruct ... bind request params to struct
func BindRequestParamsToStruct(outputStruct interface{}, params url.Values, tag string) error {
	// 1. If params empty, do nothing
	if len(params) == 0 {
		return nil
	}

	// 2. Reflect outputStruct, reject if it is not Struct
	rt := reflect.TypeOf(outputStruct).Elem()
	if rt.Kind() != reflect.Struct {
		return fmt.Errorf("not support struct")
	}

	// 3. Explore each field to set value
	rElem := reflect.ValueOf(outputStruct).Elem()
	for i := 0; i < rt.NumField(); i++ {
		rField := rt.Field(i)
		// 3.1. Case nested struct
		if rField.Type.Kind() == reflect.Struct {
			rtChild := rField.Type
			rElemChild := rElem.Addr().Elem()
			for j := 0; j < rtChild.NumField(); j++ {
				rFieldChild := rtChild.Field(j)
				if err := reflectionSetFieldValueOfStruct(rElemChild, rFieldChild, tag, params); err != nil {
					return err
				}
			}

			continue
		}

		// 3.2. Case slice/array
		if rField.Type.Kind() == reflect.Slice || rField.Type.Kind() == reflect.Array {
			rKey := strings.Split(rField.Tag.Get(tag), ",")[0] // use split to ignore tag "options" like omitempty, etc.
			values := params[rKey]
			if len(values) == 0 {
				continue
			}
			rFieldValue := rElem.FieldByName(rField.Name)

			rItemType := rField.Type.Elem()
			rItemValue := reflect.MakeSlice(rField.Type, len(values), len(values))
			for idx, value := range values {
				err := reflectionSetFieldValue(rKey, rItemType, rItemValue.Index(idx), value)
				if err != nil {
					return err
				}
			}

			rFieldValue.Set(rItemValue)
			continue
		}

		//3.3. Case field is not struct
		if err := reflectionSetFieldValueOfStruct(rElem, rField, tag, params); err != nil {
			return err
		}
	}
	return nil
}

// private method ------------------------------------------------------------------------------------------------------

func reflectionSetFieldValueOfStruct(rElem reflect.Value, rField reflect.StructField, tag string, params url.Values) error {
	rFieldValue := rElem.FieldByName(rField.Name)
	rKey := strings.Split(rField.Tag.Get(tag), ",")[0] // use split to ignore tag "options" like omitempty, etc.

	// 3.2. If key existed, set value
	if value := params.Get(rKey); len(value) > 0 {
		err := reflectionSetFieldValue(rKey, rField.Type, rFieldValue, value)
		if err != nil {
			logger.Errorf("bindRequestParamsToStruct ... error %+v", err)
			return err
		}
	}
	return nil
}

// reflectionSetFieldValue ... set value for field in struct
// support field type: String, *String, Int, *Int, Int32, *Int32, Int64, *Int64
func reflectionSetFieldValue(rKey string, rFieldType reflect.Type, rFieldValue reflect.Value, value string) error {
	// 1. Check field is pointer or not
	isPointer := false
	fieldType := rFieldType.Kind()
	if fieldType == reflect.Pointer {
		isPointer = true
		fieldType = rFieldType.Elem().Kind()
	}

	// 2. Set value for field
	switch fieldType {
	case reflect.String:
		if isPointer {
			rFieldValue.Set(reflect.ValueOf(&value))
		} else {
			rFieldValue.SetString(value)
		}
	case reflect.Int64:
		if intVal, err := strconv.Atoi(value); err == nil {
			if isPointer {
				temp := int64(intVal)
				rFieldValue.Set(reflect.ValueOf(&temp))
			} else {
				rFieldValue.SetInt(int64(intVal))
			}
		} else {
			return fmt.Errorf("%v is wrong type", rKey)
		}
	case reflect.Int32:
		if intVal, err := strconv.Atoi(value); err == nil {
			if isPointer {
				temp := int32(intVal)
				rFieldValue.Set(reflect.ValueOf(&temp))
			} else {
				rFieldValue.Set(reflect.ValueOf(int32(intVal)))
			}
		} else {
			return fmt.Errorf("%v is wrong type", rKey)
		}
	case reflect.Int:
		if intVal, err := strconv.Atoi(value); err == nil {
			if isPointer {
				temp := intVal
				rFieldValue.Set(reflect.ValueOf(&temp))
			} else {
				rFieldValue.Set(reflect.ValueOf(intVal))
			}
		} else {
			return fmt.Errorf("%v is wrong type", rKey)
		}
	default:
		return fmt.Errorf("%v with type %v is not support", rKey, fieldType)
	}
	return nil
}
