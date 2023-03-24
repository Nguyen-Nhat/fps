package fileservice

import (
	"bytes"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
)

type (
	// IService ...
	IService interface {
		UploadFileWithBytesData(*bytes.Buffer, string) (string, error)
	}

	// Service ...
	Service struct {
		client IClient
	}
)

var _ IService = &Service{}

// NewService ...
func NewService(client IClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) UploadFile(byteData *bytes.Buffer, fileName string) (string, error) {
	// 1. Build request
	req := uploadFileRequest{
		FileData: byteData.Bytes(),
		FileName: fileName,
	}

	// 2. Request upload
	res, err := s.client.uploadFile(req)
	if err != nil {
		logger.Errorf("Upload file to File Service failed: %v", err)
		return "", err
	}

	// 3. Return
	logger.Infof("Response = %v", res)
	return res.Url, nil
}

func (s *Service) UploadFileWithBytesData(dataByteBuffer *bytes.Buffer, resultFileName string) (string, error) {
	resultFileUrl, err := s.UploadFile(dataByteBuffer, resultFileName)
	if err != nil {
		logger.Errorf("Cannot upload file %v, got: %v", resultFileName, err)
		return "", err
	}

	return resultFileUrl, nil
}
