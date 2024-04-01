package fileprocessing

type GetFileProcessHistoryDTO struct {
	ClientID        int32
	SellerId        int32
	CreatedBy       string
	Page            int
	PageSize        int
	CreatedByEmails []string
	ProcessFileIds  []int
	SearchFileName  string
}

type CreateFileProcessingReqDTO struct {
	ClientID       int32
	FileURL        string
	DisplayName    string
	CreatedBy      string
	FileParameters string
	SellerID       int32
}

type CreateFileProcessingResDTO struct {
	ProcessFileID int32
}
