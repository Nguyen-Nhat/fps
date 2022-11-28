package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"io"
	"log"
	mr "math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

// Type ----------------------------------------------------------------------------------------------------------------

type (
	FileContent struct {
		FieldName string
		FileName  string
		Data      []byte
	}

	FileName struct {
		FullName  string
		Name      string
		Extension string
	}
)

func (f *FileName) FullNameWithSuffix(suffix string) string {
	return fmt.Sprintf("%s%s.%s", f.Name, suffix, f.Extension)
}

// ---------------------------------------------------------------------------------------------------------------------

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

func UploadFile[RES any](client *http.Client, urlPath string, content FileContent) (*RES, error) {
	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := createFormFile(writer, content.FieldName, content.FileName)
	if err != nil {
		log.Printf("make request failed: %d", err)
		return nil, err
	}
	_, _ = fw.Write(content.Data)

	// Close multipart writer.
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, err := client.Do(req)
	if err != nil {
		logger.Errorf("error request %v", err)
		return nil, err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()

	respBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		logger.Errorf("error read response body %v", err)
		return nil, err
	}
	logger.Infof(" ===== Upload file response body: %s", respBody)

	if rsp.StatusCode != http.StatusOK {
		logger.Infof("Request failed with response code: %d", rsp.StatusCode)
		return nil, fmt.Errorf("request failed with http status %v", rsp.StatusCode)
	}

	// 6. Convert response body to Entity
	var respBodyObj RES
	if err := json.Unmarshal(respBody, &respBodyObj); err != nil {
		logger.Errorf("===== http: Decode to entity error: %+v\n", err.Error())
		return nil, err
	}
	return &respBodyObj, nil
}

func createFormFile(w *multipart.Writer, fieldName, fileName string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	return w.CreatePart(h)
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(byteLength int) ([]byte, error) {
	b := make([]byte, byteLength)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// Length of string = ceil(1.333 * byteLength)
// securely generated random string.
func GenerateRandomString(byteLength int) (string, error) {
	b, err := generateRandomBytes(byteLength)
	return base64.RawStdEncoding.EncodeToString(b), err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(numberChars int) string {
	b := make([]byte, numberChars)
	for i := range b {
		b[i] = letterBytes[mr.Intn(len(letterBytes))]
	}
	return string(b)
}

func ExtractFileName(filePath string) FileName {
	matches := constant.FileNameRegex.FindStringSubmatch(filePath)
	if len(matches) >= 4 {
		return FileName{
			FullName:  urlDecoded(matches[1]),
			Name:      urlDecoded(matches[2]),
			Extension: matches[3],
		}
	}
	return FileName{"unknown", "unknown", "unknown"}
}

func urlDecoded(filePath string) string {
	result, err := url.QueryUnescape(filePath)
	if err != nil {
		result = filePath
	}
	return result
}
