package fileservice

import (
	config "git.teko.vn/loyalty-system/loyalty-file-processing/configs"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"testing"
)

// For temporary check ... removeThisPrefixForRunningThisTest_
func TestService_LoadFileAwardPoint(t *testing.T) {
	cfg := config.FileServiceConfig{
		Endpoint: "https://files.dev.tekoapis.net",
		Paths: config.FileServicePaths{
			UploadDoc: "/upload/doc",
		},
	}
	client := NewClient(cfg)
	service := NewService(client)

	fileUrl := "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/9/Fumart%20Loyalty%20-%20nap%20diem%20KH.xlsx"
	sheet, err := service.LoadFileAwardPoint(fileUrl)
	if err != nil {
		logger.Errorf("Error %v", err)
	} else {
		logger.Infof("Sheet: %v", sheet)
	}
}

func TestService_UploadFileAwardPointResult(t *testing.T) {
	cfg := config.FileServiceConfig{
		Endpoint: "https://files.dev.tekoapis.net",
		Paths: config.FileServicePaths{
			UploadDoc: "/upload/doc",
		},
	}
	client := NewClient(cfg)
	service := NewService(client)

	fileUrl := "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/9/Fumart%20Loyalty%20-%20nap%20diem%20KH_6.xlsx"

	fabResults := []dto.FileAwardPointResultRow{
		{RowId: 1, Phone: "0393227489", Point: 1, Note: "This is note....", Error: "This is error ..."},
		{RowId: 2, Phone: "0393227433", Point: 32, Note: "This is note 23 ....", Error: "This is error 23 ..."},
		{RowId: 3, Phone: "0393227433", Point: 32, Note: "This is note 23 ....", Error: "This is error 23 ..."},
		{RowId: 4, Phone: "0393227433", Point: 32, Note: "This is note 23 ....", Error: "This is error 23 ..."},
		{RowId: 5, Phone: "0393227433", Point: 32, Note: "This is note 23 ....", Error: "This is error 23 ..."},
		{RowId: 6, Phone: "0393227433", Point: 32, Note: "This is note 23 ....", Error: "This is error 23 ..."},
		{RowId: 7, Phone: "0393227433", Point: 32, Note: "This is note 23 ....", Error: "This is error 23 ..."},
		{RowId: 8, Phone: "0393227433", Point: 32, Note: "This is note 23 ....", Error: "This is error 23 ..."},
	}

	url, err := service.AppendErrorAndUploadFileAwardPointResult(fabResults, fileUrl)
	if err != nil {
		logger.Errorf("Error %v", err)
	} else {
		logger.Infof("Url: %v", url)
	}
}
