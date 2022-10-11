package fileservice

type uploadFileRequest struct {
	FileData []byte
	FileName string
}

type uploadFileResponse struct {
	Url string `json:"url"`
}
