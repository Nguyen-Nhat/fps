package customFunc

// Response ------------------------------------------------------------------------------------------------------------

type SiteInfo struct {
	Id             int    `json:"id"`
	SellerSiteCode string `json:"sellerSiteCode"`
}

type GetSiteResponse struct {
	Data []SiteInfo `json:"data"`
}
