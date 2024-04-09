package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mr "math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	t "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customtype"
	"moul.io/http2curl"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

const (
	isPrivateStr   = "isPrivate"
	isPrivateValue = "true"
	suffixResult   = "_result"
)

const (
	XlsxContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	CsvContentType  = "text/csv"
)

// Type ----------------------------------------------------------------------------------------------------------------

type (
	FileContent struct {
		FieldName   string
		FileName    string
		Data        []byte
		ContentType string
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
	header map[string]string,
	requestParams map[string]string,
	requestBody *REQ,
) (int, *RES, error) {
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
		return 0, nil, err
	}

	// 3. Set Header
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	// 4. Set request params
	if len(requestParams) > 0 {
		query := req.URL.Query()
		for paramField, paramValue := range requestParams {
			query.Add(paramField, paramValue)
		}
		req.URL.RawQuery = query.Encode()
	}

	curl := getCurlCommand(req)
	logger.Infof("=====> curl %+v\n", curl) // for debugging, todo remove

	// 4. Send request
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("===== http: send request error: %+v\n", err.Error())
		return 0, nil, err
	}

	// 5. Ready response body
	defer func() {
		_ = resp.Body.Close()
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("===== http: ready response body error: %+v\n", err.Error())
		return resp.StatusCode, nil, err
	}

	// 6. Convert response body to Entity
	var respBodyObj RES
	if err := json.Unmarshal(respBody, &respBodyObj); err != nil {
		logger.Errorf("===== http: Decode to entity error: %+v\n", err.Error())
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, &respBodyObj, nil
}

// SendHTTPRequestRaw ... return (httpStatusCode, responseBody, curlCommand, error)
func SendHTTPRequestRaw(
	client *http.Client,
	method, url string,
	header map[string]interface{},
	requestParams []t.Pair[string, string],
	requestBody map[string]interface{},
) (int, string, string, error) {
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
		return 0, "", "", err
	}

	// 3. Set Header
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	// 4. Set request params
	if len(requestParams) > 0 {
		query := req.URL.Query()
		for _, param := range requestParams {
			query.Add(param.Key, param.Value)
		}
		req.URL.RawQuery = query.Encode()
	}

	// 5. Send request
	curl := getCurlCommand(req)

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("===== http: send request error: %+v\n", err.Error())
		return 0, "", curl, err
	}

	// 6. Ready response body
	defer func() {
		_ = resp.Body.Close()
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("===== http: ready response body error: %+v\n", err.Error())
		return resp.StatusCode, "", curl, err
	}

	return resp.StatusCode, string(respBody), curl, nil
}

func UploadFile[RES any](client *http.Client, urlPath string, content FileContent) (*RES, error) {
	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.WriteField(isPrivateStr, isPrivateValue)
	if err != nil {
		log.Printf("make request failed: %d", err)
		return nil, err
	}
	fw, err := createFormFile(writer, content.FieldName, content.FileName, content.ContentType)
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

func createFormFile(w *multipart.Writer, fieldName, fileName string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", contentType)
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
	// get file name from private url
	fileName, err := extractFileNameFromPrivateUrl(filePath)
	if err == nil {
		return fileName
	}
	return FileName{"unknown", "unknown", "unknown"}
}

func extractFileNameFromPrivateUrl(url string) (f FileName, err error) {
	r, err := http.Get(url)
	if err != nil {
		return f, err
	}
	defer func() {
		_ = r.Body.Close()
	}()

	// Get content from Content-Disposition
	contentDisposition := r.Header.Get("Content-Disposition")
	if contentDisposition == constant.EmptyString {
		return f, nil
	}

	// Parse Content-Disposition
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return f, err
	}

	// Get file name
	f.FullName = urlDecoded(params["filename"])
	items := strings.Split(f.FullName, constant.SplitByDot)
	if len(items) == 0 {
		return f, nil
	}
	f.Extension = items[len(items)-1]
	f.Name = strings.Join(items[:len(items)-1], constant.SplitByDot)
	return f, nil
}

func urlDecoded(filePath string) string {
	result, err := url.QueryUnescape(filePath)
	if err != nil {
		result = filePath
	}
	return result
}

func JsonString[T any](data T) string {
	out, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Cannot marshal struct to string, err=%v", err)
		return fmt.Sprintf("error: %v", err)
	}

	return string(out)
}

func BatchExecuting[T any](batchSize int, listData []T, execute func([]T) error) error {
	for head := 0; head < len(listData); head += batchSize {
		tail := head + batchSize
		if tail > len(listData) {
			tail = len(listData)
		}

		batch := listData[head:tail]

		logger.Infof("Execute %v item (%v -> %v)\n", len(batch), head, tail)
		err := execute(batch)

		if err != nil {
			return err
		}

	}

	return nil
}

func BatchExecutingReturn[T any, R any](batchSize int, listData []T, execute func([]T, ...interface{}) ([]R, error), extraData ...interface{}) ([]R, error) {
	var result []R
	for head := 0; head < len(listData); head += batchSize {
		tail := head + batchSize
		if tail > len(listData) {
			tail = len(listData)
		}

		batch := listData[head:tail]

		logger.Infof("Execute %v item (%v -> %v)\n", len(batch), head, tail)

		res, err := execute(batch, extraData...)
		if err != nil {
			return nil, err
		}

		result = append(result, res...)
	}

	return result, nil
}

func getCurlCommand(req *http.Request) string {
	curlCommand, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return err.Error()
	}

	return curlCommand.String()
}

func CloneMap[D any](root map[string]D) map[string]D {
	targetMap := make(map[string]D)
	for key, value := range root {
		targetMap[key] = value
	}
	return targetMap
}

func CloneArray[D any](root []D) []D {
	return append([]D{}, root...)
}

func TrimSpaceAndToLower(input string) string {
	out := strings.TrimSpace(input)
	out = strings.ToLower(out)
	return out
}

func EqualsIgnoreCase(s1 string, s2 string) bool {
	return TrimSpaceAndToLower(s1) == TrimSpaceAndToLower(s2)
}

func GetResultFileName(fileName string) string {
	if len(fileName) == 0 {
		return constant.EmptyString
	}
	fileNameExtract := strings.Split(fileName, constant.SplitByDot)
	fileNameExtract[0] += suffixResult
	return strings.Join(fileNameExtract, constant.SplitByDot)
}

func Contains[T comparable](arr []T, val T) bool {
	for _, e := range arr {
		if e == val {
			return true
		}
	}

	return false
}
