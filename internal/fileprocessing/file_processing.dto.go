package fileprocessing

type CreateFileProcessingReqDTO struct {
	ClientID    int64
	FileURL     string
	DisplayName string
	CreatedBy   string
}

type CreateFileProcessingResDTO struct {
	ProcessFileID int32
}
