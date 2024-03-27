package customFunc

// Response ------------------------------------------------------------------------------------------------------------

type SupplierInfo struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
}

type GetSupplierResponse struct {
	Data struct {
		Suppliers []SupplierInfo `json:"suppliers"`
	} `json:"data"`
}
