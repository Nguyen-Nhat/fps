package converter

import (
	"encoding/json"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

func StringToMap(dataName string, dataRaw string, allowEmpty bool) (map[string]string, error) {
	if allowEmpty && len(dataRaw) <= 0 {
		return map[string]string{}, nil
	}

	headerMap := map[string]string{}
	if err := json.Unmarshal([]byte(dataRaw), &headerMap); err != nil {
		logger.Errorf("error when convert %v: value=%v, err=%v", dataName, dataRaw, err)
		return nil, fmt.Errorf("cannot convert %v", dataName)
	}
	return headerMap, nil
}

func StringJsonToStruct[OUT any](dataName string, dataRaw string, data OUT) (*OUT, error) {
	if err := json.Unmarshal([]byte(dataRaw), &data); err != nil {
		logger.Errorf("error when convert %v: value=%v, err=%v", dataName, dataRaw, err)
		return &data, fmt.Errorf("cannot convert %v", dataName)
	}
	return &data, nil
}
