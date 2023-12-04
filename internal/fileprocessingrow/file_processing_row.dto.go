package fileprocessingrow

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/request"
)

type CreateProcessingFileRowJob struct {
	FileID       int
	RowIndex     int
	RowDataRaw   string
	TaskIndex    int
	TaskMapping  string
	GroupByValue string
}

type UpdateAfterExecutingByJob struct {
	Task ProcessingFileRow // original data

	TaskMapping  string
	RequestCurl  string
	ResponseRaw  string
	Status       int16
	ErrorDisplay string
	ExecutedTime int64 // unit milliseconds
}

type StatisticData struct {
	IsFinished    bool
	ErrorDisplays map[int]string

	TotalProcessed int
	TotalSuccess   int
	TotalFailed    int
	TotalWaiting   int
}

type GetListFileRowsRequest struct {
	request.PageRequest
}

type GetListFileRowsItem struct {
	FileID       int             `json:"fileID"`
	RowIndex     int             `json:"rowIndex"`
	RowDataRaw   string          `json:"rowDataRaw"`
	ExecutedTime int             `json:"executedTime"`
	Tasks        []TaskInRowItem `json:"tasks"`
}

type TaskInRowItem struct {
	TaskIndex       int    `json:"taskIndex"`
	TaskRequestCurl string `json:"taskRequestCurl"`
	TaskResponseRaw string `json:"taskResponseRaw"`
	TaskName        string `json:"taskName"`
	GroupByValue    string `json:"groupByValue"`
	Status          int16  `json:"status"`
	ErrorDisplay    string `json:"errorDisplay"`
	ExecutedTime    int    `json:"executedTime"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
}
