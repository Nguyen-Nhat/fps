package fileprocessing

type CreateFileProcessingReqDTO struct {
	ClientID    int32
	FileURL     string
	DisplayName string
	CreatedBy   string
}

type CreateFileProcessingResDTO struct {
	ProcessFileID int32
}
