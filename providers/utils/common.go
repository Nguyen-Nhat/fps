package utils

import (
	"bytes"
	"encoding/json"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"io"
	"net/http"
)

func HiddenString(input string, numberOfTailChar int) string {
	if numberOfTailChar >= len(input) {
		return input
	}
	return "***" + input[len(input)-numberOfTailChar:]
}

func SendHTTPRequest[REQ any, RES any](
	client *http.Client,
	method, url string,
	header map[string]string, requestBody *REQ,
) (*RES, error) {
	// 1. Build body
	var bodyIO *bytes.Buffer
	if requestBody == nil {
		bodyIO = bytes.NewBuffer([]byte{})
	} else {
		requestBytes, _ := json.Marshal(requestBody)
		bodyIO = bytes.NewBuffer(requestBytes)
	}

	// 2. Build request
	req, err := http.NewRequest(method, url, bodyIO)
	if err != nil {
		return nil, err
	}

	// 3. Set Header
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	// 4. Send request
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("===== http: send request error: %+v\n", err.Error())
		return nil, err
	}

	// 5. Ready response body
	defer func() {
		_ = resp.Body.Close()
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("===== http: ready response body error: %+v\n", err.Error())
		return nil, err
	}

	// 6. Convert response body to Entity
	var respBodyObj RES
	if err := json.Unmarshal(respBody, &respBodyObj); err != nil {
		logger.Errorf("===== http: Decode to entity error: %+v\n", err.Error())
		return nil, err
	}
	return &respBodyObj, nil
}
