package response

import "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/request"

type PaginationNew struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func GetPaginationNew(total int, req request.PageRequest) PaginationNew {
	return PaginationNew{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
}
