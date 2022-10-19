package fileawardpoint

import (
	"errors"
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/response"
	"net/http"
)

// Request DTO =========================================================================================================

// GetFileAwardPointDetailRequest ...
type GetFileAwardPointDetailRequest struct {
	Id int `json:"id"`
}

func (a *GetFileAwardPointDetailRequest) Bind(_ *http.Request) error {
	// validate Id missing
	if a.Id == 0 {
		return errors.New("required id")
	}

	return nil
}

// Response DTO ========================================================================================================

// GetFileAwardPointDetailResponse is the response payload for the User data model.
//
// In the userResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type GetFileAwardPointDetailResponse struct {
	Id                int    `json:"id"`
	MerchantId        int64  `json:"merchantId"`
	DisplayName       string `json:"displayName"`
	FileUrl           string `json:"fileUrl"`
	ResultFileUrl     string `json:"resultFileUrl"`
	Status            string `json:"status"`
	StatsTotalRow     int32  `json:"statsTotalRow"`
	StatsTotalSuccess int32  `json:"StatsTotalSuccess"`
	CreatedAt         int64  `json:"createdAt"`
	CreatedBy         string `json:"createdBy"`
}

func (ur *GetFileAwardPointDetailResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

type GetListFileAwardPointData struct {
	FileAwardPoints []FileAwardPoint    `json:"fileAwardPoints"`
	Pagination      response.Pagination `json:"pagination"`
}

type FileAwardPoint struct {
	MerchantId        int64  `json:"merchantId"`
	FileAwardPointId  int    `json:"fileAwardPointId"`
	FileDisplayName   string `json:"fileDisplayName"`
	FileUrl           string `json:"fileUrl"`
	ResultFileUrl     string `json:"resultFileUrl"`
	Status            string `json:"status"`
	StatsTotalRow     int32  `json:"statsTotalRow"`
	StatsTotalSuccess int32  `json:"statsTotalSuccess"`
	CreatedAt         int64  `json:"createdAt"`
	CreatedBy         string `json:"createdBy"`
}

// CreateFileAwardPointDetailRequest ...
type CreateFileAwardPointDetailRequest struct {
	MerchantID int64  `json:"merchantId"`
	FileUrl    string `json:"fileUrl"`
	Note       string `json:"note"`
}

func (c *CreateFileAwardPointDetailRequest) Bind(_ *http.Request) error {
	// validate Id missing
	if c.MerchantID == 0 {
		return ErrMerchantIDRequired
	}

	// validate fileUrl missing
	if c.FileUrl == "" {
		return ErrFileUrlRequired
	}

	// validate Note length
	if len(c.Note) > 255 {
		return ErrNoteOverMaxLength
	}

	return nil
}

// CreateFileAwardPointDetailResponse ...
type CreateFileAwardPointDetailResponse struct {
	FileAwardPointID int `json:"fileAwardPointId"`
}

var (
	ErrMerchantIDRequired = errors.New("merchant id is required")
	ErrFileUrlRequired    = errors.New("file url is required")
	ErrNoteOverMaxLength  = errors.New("note over 255 character")
)

const (
	FapStatusInit       = "INIT"
	FapStatusProcessing = "PROCESSING"
	FapStatusFailed     = "FAILED"
	FapStatusFinished   = "FINISHED"
)
