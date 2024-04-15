package csv

import (
	"encoding/csv"
	"net/http"
)

func LoadCSVByURL(fileURL string) ([][]string, error) {
	r, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Body.Close()
	}()
	return csv.NewReader(r.Body).ReadAll()
}
