package funcClient10

// input ---------------------------------------------------------------------------------------------------------------

type ItemInput struct {
	SellerSku string  `json:"sellerSku"`
	UomName   string  `json:"uomName"`
	Quantity  float64 `json:"quantity"`
}

// Output --------------------------------------------------------------------------------------------------------------

type ItemOutput struct {
	Sku      string  `json:"sku"`
	Quantity float64 `json:"quantity"`
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
	SellerId  int    `json:"sellerId"`
}
