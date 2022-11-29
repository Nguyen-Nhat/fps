package response

import "math"

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPage   int `json:"totalPage"`
}

func GetPagination(total int, page int, size int) *Pagination {
	return &Pagination{
		CurrentPage: page,
		PageSize:    size,
		TotalItems:  total,
		TotalPage:   int(math.Ceil(float64(total) / float64(size))),
	}
}
