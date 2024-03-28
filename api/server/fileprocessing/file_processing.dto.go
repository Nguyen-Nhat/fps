package fileprocessing

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// Request DTO =========================================================================================================

// CreateFileProcessingRequest ...
type CreateFileProcessingRequest struct {
	ClientID        int32  `json:"clientId"`
	FileURL         string `json:"fileUrl"`
	FileDisplayName string `json:"fileDisplayName"`
	CreatedBy       string `json:"createdBy"`
	Parameters      string `json:"parameters"`
	SellerID        int32  `json:"sellerId"`
}

// Take from https://stackoverflow.com/a/36922225
func isJSONString(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func (c *CreateFileProcessingRequest) Bind(r *http.Request) error {
	logger.Infof("===== Request CreateFileProcessingRequest: \n%+v\n", *c)

	// validate Id missing
	if c.ClientID == 0 {
		return ErrClientIDRequired
	}

	// validate fileUrl missing
	if c.FileURL == "" {
		return ErrFileUrlRequired
	}
	// validate createdBy missing
	if c.CreatedBy == "" {
		return ErrCreatedByRequired
	}

	// validate fileUrl length
	if len(c.FileURL) > FileUrlMaxLength {
		return ErrFileUrlOverMaxLength
	}

	// validate displayName length
	if len(c.FileDisplayName) > DisplayNameMaxLength {
		return ErrDisplayNameOverMaxLength
	}
	// validate createdBy length
	if len(c.CreatedBy) > CreatedByMaxLength {
		return ErrCreatedByOverMaxLength
	}
	// validate parameters field, valid when it is "" or is in JSON format
	if c.Parameters != "" && !isJSONString(c.Parameters) {
		return ErrParametersIsNotJson
	}
	return nil
}

type GetFileProcessHistoryRequest struct {
	ClientID        int32    `json:"clientId"`
	SellerID        int32    `json:"sellerId"`
	CreatedBy       string   `json:"createdBy"`
	Page            int      `json:"page"`
	PageSize        int      `json:"size"`
	CreatedByEmails []string `json:"createdByEmails"`
	ProcessFileIds  []int    `json:"processFileIds"`
	SearchFileName  string   `json:"searchFileName"`
}

func bindAndValidateRequestParams(r *http.Request, data *GetFileProcessHistoryRequest) error {
	// 1.Bind request params to struct
	params := r.URL.Query()
	if err := utils.BindRequestParamsToStruct(data, params, "json"); err != nil {
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return error2.ErrInvalidRequestWithError(err)
	}
	data.InitDefaultValue()

	// 2.Validate request params
	if data.ClientID == 0 {
		err := fmt.Errorf("missing required param: clientId")
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return error2.ErrInvalidRequestWithError(err)
	}

	if data.Page < 0 || data.Page > constant.PaginationMaxPage {
		err := fmt.Errorf("request field is out of range: page")
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return error2.ErrInvalidRequestWithError(err)
	}

	if data.PageSize < 0 || data.PageSize > constant.PaginationMaxSize {
		err := fmt.Errorf("request field is out of range: size")
		logger.Errorf("bindAndValidateRequestParams ... error %+v", err)
		return error2.ErrInvalidRequestWithError(err)
	}
	return nil
}

func (c *GetFileProcessHistoryRequest) InitDefaultValue() {
	if c.Page == 0 {
		c.Page = 1
	}

	if c.PageSize == 0 {
		c.PageSize = 10
	}
}

// Response DTO ========================================================================================================

// CreateFileProcessingResponse ...
type CreateFileProcessingResponse struct {
	ProcessFileID int64 `json:"processFileId"`
}

type GetFileProcessHistoryData struct {
	ProcessingFiles []ProcessingHistoryFile `json:"processingFiles"`
	Pagination      response.Pagination     `json:"pagination"`
}

type ProcessingHistoryFile struct {
	ClientId            int32  `json:"clientId"`
	ProcessingFileId    int    `json:"processingFileId"`
	FileDisplayName     string `json:"fileDisplayName"`
	FileUrl             string `json:"fileUrl"`
	ResultFileUrl       string `json:"resultFileUrl"`
	Status              string `json:"status"`
	SellerID            int32  `json:"sellerId"`
	StatsTotalRow       int32  `json:"statsTotalRow"`
	StatsTotalProcessed int32  `json:"statsTotalProcessed"`
	StatsTotalSuccess   int32  `json:"statsTotalSuccess"`
	ErrorDisplay        string `json:"errorDisplay"`
	CreatedAt           int64  `json:"createdAt"`
	CreatedBy           string `json:"createdBy"`
	FinishedAt          int64  `json:"finishedAt"`
}

// Error ===============================================================================================================

var (
	ErrClientIDRequired         = errors.New("clientId is required")
	ErrFileUrlRequired          = errors.New("fileUrl is required")
	ErrCreatedByRequired        = errors.New("createdBy is required")
	ErrFileUrlOverMaxLength     = fmt.Errorf("fileUrl over %d character", FileUrlMaxLength)
	ErrDisplayNameOverMaxLength = fmt.Errorf("fileDisplayName over %d character", DisplayNameMaxLength)
	ErrCreatedByOverMaxLength   = fmt.Errorf("createdBy over %d characters", CreatedByMaxLength)
	ErrParametersIsNotJson      = errors.New("parameters field isn't in json format")
)

var (
	FileUrlMaxLength     = 255
	DisplayNameMaxLength = 255
	CreatedByMaxLength   = 255
)

const (
	FpStatusInit       = "INIT"
	FpStatusProcessing = "PROCESSING"
	FpStatusFailed     = "FAILED"
	FpStatusFinished   = "FINISHED"
)
