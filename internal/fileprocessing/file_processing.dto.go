package fileprocessing

type GetFileProcessHistoryDTO struct {
	ClientID  int32
	CreatedBy string
	Page      int
	PageSize  int
}

type CreateFileProcessingReqDTO struct {
	ClientID       int32
	FileURL        string
	DisplayName    string
	CreatedBy      string
	FileParameters string
}

type CreateFileProcessingResDTO struct {
	ProcessFileID int32
}
