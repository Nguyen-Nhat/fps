package fileprocessing

import "golang.org/x/text/language"

type GetFileProcessHistoryDTO struct {
	ClientID        int32 // deprecated
	ClientIds       []int32
	SellerId        int32
	CreatedBy       string
	Page            int
	PageSize        int
	CreatedByEmails []string
	ProcessFileIds  []int
	SearchFileName  string
	MerchantId      string
	TenantId        string
}

type CreateFileProcessingReqDTO struct {
	ClientID       int32
	FileURL        string
	DisplayName    string
	CreatedBy      string
	FileParameters string
	SellerID       int32
	MerchantId     string
	TenantId       string
	AcceptLanguage language.Tag
}

type CreateFileProcessingResDTO struct {
	ProcessFileID int32
}
