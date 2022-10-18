package fileservice

import (
	"bytes"
	"fmt"
	"strconv"

	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"github.com/xuri/excelize/v2"
)

type (
	// IService ...
	IService interface {
		UploadFileAwardPointResult([]dto.FileAwardPointRow, string) (string, error)
		LoadFileAwardPoint(string) (*dto.Sheet[dto.FileAwardPointRow], error)
		UploadFileAwardPointError([]dto.ErrorRow, string) (string, error)
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

func (s *Service) UploadFileAwardPointResult(fabResults []dto.FileAwardPointRow, previousResultFileUrl string) (string, error) {
	// 1. Load previous result file
	var previousFAPResults []dto.FileAwardPointRow

	// 2. Mix previous result with current result
	// todo need to check this case, we can append data (as below line), override data or mix data
	combinedFAPResults := fabResults
	if previousResultFileUrl != "" {
		combinedFAPResults = append(combinedFAPResults, previousFAPResults...)
	}

	// 3. Convert data to bytes
	dataByteBuffer, err := makeResultFile(combinedFAPResults)
	if err != nil {
		logger.Errorf("Failed to convert data to excel %v", err)
	}

	// 4. Build request
	fileName := utils.ExtractFileName(previousResultFileUrl)
	fileUrl, err := s.UploadFile(dataByteBuffer, fileName.FullName)
	if err != nil {
		logger.Errorf("Cannot upload file, got: %v", err)
		return "", err
	}

	return fileUrl, nil
}

func (s *Service) UploadFileAwardPointError(errorRows []dto.ErrorRow, fileName string) (string, error) {
	// 1. Convert data to bytes
	dataByteBuffer, err := makeResultFileByError(errorRows)
	if err != nil {
		logger.Errorf("Failed to convert data to excel %v", err)
	}

	// 4. Build request
	fileUrl, err := s.UploadFile(dataByteBuffer, fileName)
	if err != nil {
		logger.Errorf("Cannot upload file, got: %v", err)
		return "", err
	}

	return fileUrl, nil
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

func (s *Service) LoadFileAwardPoint(url string) (*dto.Sheet[dto.FileAwardPointRow], error) {
	// 1. Build metadata config for Excel
	dataIndexStart := 3
	fileAwardPointMT := dto.FileAwardPointMetadata{
		Phone: dto.CellData[string]{
			ColumnName: "Phone number (*)",
			Constrains: dto.Constrains{IsRequired: true},
		},
		Point: dto.CellData[int]{
			ColumnName: "Points (*)",
			Constrains: dto.Constrains{IsRequired: true},
		},
		Note: dto.CellData[string]{
			ColumnName: "Note",
			Constrains: dto.Constrains{IsRequired: false},
		},
	}

	// 2. Load file from URL -> result is [][]string
	sheetData, err := excel.LoadExcelByUrl(url)
	if err != nil {
		return nil, err
	}

	// 3. Convert data to struct
	return excel.ConvertToStruct[
		dto.FileAwardPointMetadata,
		dto.FileAwardPointRow,
		dto.Converter[dto.FileAwardPointMetadata, dto.FileAwardPointRow],
	](dataIndexStart, &fileAwardPointMT, sheetData)
}

// Private method ------------------------------------------------------------------------------------------------------

func makeResultFile(rows []dto.FileAwardPointRow) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	// Create a new sheet.
	sheetName := "Sheet1"
	index := f.NewSheet(sheetName)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Set header
	_ = f.SetSheetRow(sheetName, "A1", &[]interface{}{"STT", "Phone number (*)", "Points (*)", "Note", "Error"})
	_ = f.SetSheetRow(sheetName, "A2", &[]interface{}{"Số thứ tự", "SĐT khách hàng", "Số điểm nạp", "Ghi chú giao dịch", "Kết quả"})

	// Set each row
	dataIndexStart := 3
	for rowId, row := range rows {
		axis := fmt.Sprintf("A%v", rowId+dataIndexStart)
		pointStr := strconv.Itoa(row.Point)
		err := f.SetSheetRow(sheetName, axis, &[]interface{}{row.RowId, row.Phone, pointStr, row.Note, row.Error})
		if err != nil {
			logger.Errorf("Cannot set sheet row, got %v", err)
			return nil, err
		}
	}

	return f.WriteToBuffer()
}

func makeResultFileByError(rows []dto.ErrorRow) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	// Create a new sheet.
	sheetName := "Sheet1"
	index := f.NewSheet(sheetName)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Set header
	_ = f.SetSheetRow(sheetName, "A1", &[]interface{}{"STT", "Phone number (*)", "Points (*)", "Note", "Error"})
	_ = f.SetSheetRow(sheetName, "A2", &[]interface{}{"Số thứ tự", "SĐT khách hàng", "Số điểm nạp", "Ghi chú giao dịch", "Kết quả"})

	// Set each row
	dataIndexStart := 3
	numberOfColumn := 5

	for rowId, row := range rows {
		rawDataRow := make([]interface{}, numberOfColumn)
		rawDataRow[0] = row.RowId
		for index, rowCell := range row.RowData {
			if index <= numberOfColumn-2 {
				rawDataRow[index+1] = rowCell
			}
		}

		rawDataRow[numberOfColumn-1] = row.Reason
		axis := fmt.Sprintf("A%v", rowId+dataIndexStart)
		err := f.SetSheetRow(sheetName, axis, &rawDataRow)
		if err != nil {
			logger.Errorf("Cannot set sheet row, got %v", err)
			return nil, err
		}
	}

	return f.WriteToBuffer()
}
