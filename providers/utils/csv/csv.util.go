package csv

import (
	"encoding/csv"
	"net/http"

	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

func LoadCSVByURL(fileURL string) ([][]string, error) {
	fileURL = utils.GetInternalFileURL(fileURL)
	r, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Body.Close()
	}()
	return csv.NewReader(r.Body).ReadAll()
}
