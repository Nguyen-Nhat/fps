package funcClient10

// input ---------------------------------------------------------------------------------------------------------------

type ItemInput struct {
	SellerSku  string  `json:"sellerSku"`
	UomName    string  `json:"uomName"`
	RequestQty float64 `json:"requestQty"`
}

// Output --------------------------------------------------------------------------------------------------------------

type ItemOutput struct {
	Sku        string  `json:"sku"`
	RequestQty float64 `json:"requestQty"`
}

// Response ------------------------------------------------------------------------------------------------------------

type GetSkuResponse struct {
	Data GetSkuResponseData `json:"data"`
}

type GetSkuResponseData struct {
	Products []Product `json:"products"`
}

type Product struct {
	Id        int    `json:"id"`
	Sku       string `json:"sku"`
	SellerSku string `json:"sellerSku"`
	UomName   string `json:"uomName"`
}
