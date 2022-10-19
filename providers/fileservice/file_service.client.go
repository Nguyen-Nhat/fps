package fileservice

import (
	"crypto/tls"
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"net/http"
	"time"
)

type (
	// IClient ...
	IClient interface {
		uploadFile(request uploadFileRequest) (uploadFileResponse, error)
	}

	// Client ...
	Client struct {
		conf   config.FileServiceConfig
		client *http.Client
	}
)

var _ IClient = &Client{}

// NewClient ...
func NewClient(conf config.FileServiceConfig) *Client {
	if len(conf.Endpoint) == 0 {
		panic("===== File Service Endpoint Must Not Empty")
	}

	transportCfg := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &Client{
		conf: conf,
		client: &http.Client{
			Timeout:   2 * time.Minute,
			Transport: transportCfg,
		},
	}
}

func (c *Client) uploadFile(req uploadFileRequest) (uploadFileResponse, error) {
	// 1. Build url & header
	url := c.conf.Endpoint + c.conf.Paths.UploadDoc

	// 2. Build request body
	fileContent := utils.FileContent{
		FieldName: "file",
		FileName:  req.FileName,
		Data:      req.FileData,
	}

	// 3. Send http request
	httpResp, err := utils.UploadFile[uploadFileResponse](c.client, url, fileContent)
	if err != nil {
		logger.Errorf("===== fileservice.GetTransactionByID: Call File Service Error: %s", err.Error())
		return uploadFileResponse{}, err
	}

	// 5. Response
	return *httpResp, nil
}
