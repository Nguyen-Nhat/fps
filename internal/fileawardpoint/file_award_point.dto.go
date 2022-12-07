package fileawardpoint

import "git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/pagination"

// GetFileAwardPointDetailDTO ...
type GetFileAwardPointDetailDTO struct {
	Id int
}

type GetListFileAwardPointDTO struct {
	MerchantId int
	pagination.PaginatingRequest
}

type CreateFileAwardPointReqDTO struct {
	MerchantID  int64
	FileUrl     string
	Note        string
	FileName    string
	CreatedUser string
}

type CreateFileAwardPointResDTO struct {
	FileAwardPointId int32
}
