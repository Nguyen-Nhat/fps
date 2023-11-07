package funcClient9

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

const FuncReUploadFile = "reUploadFile"

const (
	googleDriveUrl = "drive.google.com"

	// fileServiceUrl ... path is `/upload/image`
	// api-doc: https://apidoc.teko.vn/project-doc/approved/core_logic_layer/file_service_retail/version/latest/operations/post_uploads
	fileServiceUrl = "http://files-core-api.files-service/upload/image" // for calling internal service -> should move to env config
	//fileServiceUrl = "https://files.dev.tekoapis.net/upload/image" // for testing local
)

var listDomainNoNeedToReUpload = []string{"lh3.googleusercontent.com"}

var errDefault = customFunc.FuncResult{ErrorMessage: "xảy ra lỗi với đường dẫn ảnh"}

type uploadFileResponse struct {
	Url      string `json:"url"`
	ImageUrl string `json:"image_url"`
}

// ReUploadFile ...
func ReUploadFile(fullURLFile string) customFunc.FuncResult {
	// 1. Check
	if noNeedToReUpload(fullURLFile) {
		return customFunc.FuncResult{Result: fullURLFile}
	}

	// 2. Download file
	fileName, fileData, err := downloadFileStoreInMemoryBuffer(fullURLFile)
	if err != nil || len(fileData) == 0 {
		logger.Errorf("reUploadFile: downloadFile got error = %v", err)
		return errDefault
	}

	// 3. Upload file to File Service
	fileRes, err := uploadFile(fileData, fileName)
	if err != nil {
		logger.Errorf("reUploadFile: uploadFile got error = %v", err)
		return errDefault
	}

	return customFunc.FuncResult{Result: fileRes.ImageUrl}
}

// downloadFileStoreInMemoryBuffer ...
// This code doesn't download file and save to local file
// It downloads file and use memory buffer to store file
// Note: it can use too much RAM -> OOM service
func downloadFileStoreInMemoryBuffer(fullURLFile string) (string, []byte, error) {
	// 1. Init client
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
		Timeout: 10 * time.Second,
	}

	// 2. Get file url and file name
	segments := strings.Split(fullURLFile, "/")
	fileName := segments[len(segments)-1]
	fixedUrl := fullURLFile
	if strings.Contains(fullURLFile, googleDriveUrl) {
		if len(segments) >= 8 {
			fileID := segments[7]
			fixedUrl = fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", fileID)
			fileName = fmt.Sprintf("%v.jpg", fileID)
		} else {
			logger.Errorf("Invalid url: %v", fullURLFile)
			return "", nil, fmt.Errorf("đường dẫn lỗi")
		}
	}

	// 3. Get file
	resp, err := client.Get(fixedUrl)
	if err != nil {
		return "", nil, err
	} else if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("error when download file, resp = %+v", resp.Status)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Errorf("defer body.close got error %v", err)
		}
	}(resp.Body)

	// 4. Read the response body into a byte slice
	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	return fileName, fileData, nil
}

// uploadFile ...
func uploadFile(fileData []byte, fileName string) (uploadFileResponse, error) {
	// 1. Build url & header
	params := url.Values{}
	params.Add("fileName", fileName)
	params.Add("cloud", "true")

	requestUrl := fmt.Sprintf("%s?%s", fileServiceUrl, params.Encode())

	// 2. Build request body
	fileContent := utils.FileContent{
		FieldName:   "file",
		FileName:    fileName,
		Data:        fileData,
		ContentType: "image/jpg", // hardcode, need check file extension to focus to correct content type
	}

	// 3. Send http request
	client := initHttpClient()
	httpResp, err := utils.UploadFile[uploadFileResponse](client, requestUrl, fileContent)
	if err != nil {
		logger.Errorf("===== reUploadFile: Call File Service Error: %s", err.Error())
		return uploadFileResponse{}, err
	}

	// 5. Response
	return *httpResp, nil
}

// initHttpClient...
func initHttpClient() *http.Client {
	transportCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Timeout:   20 * time.Second,
		Transport: transportCfg,
	}
	return client
}

// noNeedToReUpload ...
func noNeedToReUpload(url string) bool {
	for _, domain := range listDomainNoNeedToReUpload {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return false
}
