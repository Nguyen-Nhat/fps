package fileawardpoint

import "time"

// FileAwardPoint is model of table `file_award_point`
type FileAwardPoint struct {
	Id                int       `json:"id"`
	MerchantId        int       `json:"merchant_id"`
	DisplayName       string    `json:"display_name"`
	FileUrl           string    `json:"file_url"`
	ResultFileUrl     string    `json:"result_file_url"`
	Status            int       `json:"status"`
	StatsTotalRow     int       `json:"stats_total_row"`
	StatsTotalSuccess int       `json:"stats_total_success"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         string    `json:"created_by"`
	UpdatedBy         string    `json:"updated_by"`
}

// Status ENUM ...
const (
	statusInit       = 0
	statusProcessing = 1
	statusSuccess    = 2
	statusFailed     = 3
)
