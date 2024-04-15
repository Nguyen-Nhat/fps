package filewriter

import (
	"io"
	"net/http"
)

func loadDataFromURL(fileURL string) ([]byte, error) {
	httpClient := &http.Client{}
	r, err := httpClient.Get(fileURL)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = r.Body.Close()
	}()

	return io.ReadAll(r.Body)
}
