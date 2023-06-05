package response

import "math"

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPage   int `json:"totalPage"`
}

func GetPagination(total int, page int, pageSize int) *Pagination {
	return &Pagination{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  total,
		TotalPage:   int(math.Ceil(float64(total) / float64(pageSize))),
	}
}
