package fileprocessing

import "git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/pagination"

type GetFileProcessHistoryDTO struct {
	ClientId int32
	pagination.PaginatingRequest
}

type CreateFileProcessingReqDTO struct {
	ClientID    int32
	FileURL     string
	DisplayName string
	CreatedBy   string
}

type CreateFileProcessingResDTO struct {
	ProcessFileID int32
}
