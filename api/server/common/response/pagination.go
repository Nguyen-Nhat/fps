package response

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPage   int `json:"totalPage"`
}
