package fileprocessing

import (
	"errors"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"net/http"
)

import "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"

// Request DTO =========================================================================================================

// CreateFileProcessingRequest ...
type CreateFileProcessingRequest struct {
	ClientID        int32  `json:"clientId"`
	FileURL         string `json:"fileUrl"`
	FileDisplayName string `json:"fileDisplayName"`
	CreatedBy       string `json:"createdBy"`
}

func (c *CreateFileProcessingRequest) Bind(r *http.Request) error {
	logger.Infof("===== Request CreateFileProcessingRequest: \n%v\n", utils.JsonString(r))

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

	return nil
}

// Response DTO ========================================================================================================

// CreateFileProcessingResponse ...
type CreateFileProcessingResponse struct {
	ProcessFileID int64 `json:"processFileId"`
}

// Error ===============================================================================================================

var (
	ErrClientIDRequired         = errors.New("client id is required")
	ErrFileUrlRequired          = errors.New("file url is required")
	ErrCreatedByRequired        = errors.New("created by is required")
	ErrFileUrlOverMaxLength     = fmt.Errorf("file url over %d character", FileUrlMaxLength)
	ErrDisplayNameOverMaxLength = fmt.Errorf("display name over %d character", DisplayNameMaxLength)
	ErrCreatedByOverMaxLength   = fmt.Errorf("created by over %d character", CreatedByMaxLength)
)

var (
	FileUrlMaxLength     = 255
	DisplayNameMaxLength = 255
	CreatedByMaxLength   = 255
)

type GetFileProcessHistoryData struct {
	ProcessingFiles []ProcessingHistoryFile `json:"processingFiles"`
	Pagination      response.Pagination     `json:"pagination"`
}

type ProcessingHistoryFile struct {
	ClientId          int32  `json:"clientId"`
	ProcessingFileId  int    `json:"processingFileId"`
	FileDisplayName   string `json:"fileDisplayName"`
	FileUrl           string `json:"fileUrl"`
	ResultFileUrl     string `json:"resultFileUrl"`
	Status            string `json:"status"`
	StatsTotalRow     int32  `json:"statsTotalRow"`
	StatsTotalSuccess int32  `json:"statsTotalSuccess"`
	ErrorDisplay      string `json:"errorDisplay"`
	CreatedAt         int64  `json:"createdAt"`
	CreatedBy         string `json:"createdBy"`
}

const (
	FpStatusInit       = "INIT"
	FpStatusProcessing = "PROCESSING"
	FpStatusFailed     = "FAILED"
	FpStatusFinished   = "FINISHED"
)
