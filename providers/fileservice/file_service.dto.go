package fileservice

type uploadFileRequest struct {
	FileData []byte
	FileName string
	FileType string
}

type uploadFileResponse struct {
	Url string `json:"url"`
}
