package fileawardpoint

import (
	"errors"
	"net/http"
	"time"
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
	Id                int       `json:"id"`
	MerchantId        int       `json:"merchant_id"`
	DisplayName       string    `json:"display_name"`
	FileUrl           string    `json:"file_url"`
	ResultFileUrl     string    `json:"result_file_url"`
	Status            int       `json:"status"`
	StatsTotalRow     int       `json:"stats_total_row"`
	StatsTotalSuccess int       `json:"stats_total_success"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         string    `json:"created_by"`
}

func (ur *GetFileAwardPointDetailResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
