package csv

import (
	"bytes"
	"encoding/csv"

	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

func LoadCSVByURL(fileURL string) ([][]string, error) {
	data, err := utils.GetDataFromURL(fileURL)
	if err != nil {
		return nil, err
	}

	return csv.NewReader(bytes.NewReader(data)).ReadAll()
}
