package fileservice

import (
	"bytes"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"github.com/xuri/excelize/v2"
)

type (
	// IService ...
	IService interface {
		AppendErrorAndUploadFileAwardPointResult([]dto.FileAwardPointResultRow, string) (string, error)
		LoadFileAwardPoint(string) (*dto.Sheet[dto.FileAwardPointRow], error)
		UploadFileAwardPointError([]dto.ErrorRow, string) (string, error)
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

func (s *Service) AppendErrorAndUploadFileAwardPointResult(fabResults []dto.FileAwardPointResultRow, previousResultFileUrl string) (string, error) {
	// 1. Load previous result file
	logger.Infof("fetch file %v", previousResultFileUrl)
	previousFAPResults, err := excel.LoadExcelByUrl(previousResultFileUrl)
	if err != nil {
		logger.Errorf("Failed to fetch file: ", err)
	}

	// 2. Mix data then Convert data to bytes
	dataByteBuffer, err := mixAndMakeResultFile(previousFAPResults, fabResults)
	if err != nil {
		logger.Errorf("Failed to convert data to excel %v", err)
		return "", err
	}

	// 3. Build request
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

func (s *Service) UploadFileWithBytesData(dataByteBuffer *bytes.Buffer, resultFileName string) (string, error) {
	resultFileUrl, err := s.UploadFile(dataByteBuffer, resultFileName)
	if err != nil {
		logger.Errorf("Cannot upload file %v, got: %v", resultFileName, err)
		return "", err
	}

	return resultFileUrl, nil
}

// Private method ------------------------------------------------------------------------------------------------------

func mixAndMakeResultFile(previousFAPResults [][]string, newResults []dto.FileAwardPointResultRow) (*bytes.Buffer, error) {
	// 1. Init sheet
	const dataIndexStart = 3
	exFile, sheetName := initFapResultSheetWithHeader()

	// 2. Set previous result
	if len(previousFAPResults) >= dataIndexStart-1 {
		previousFAPResults = previousFAPResults[dataIndexStart-1:] // ignore header
	} else {
		previousFAPResults = [][]string{}
	}
	for rowId, previousResultRow := range previousFAPResults {
		axis := fmt.Sprintf("A%v", rowId+dataIndexStart)
		sheetRow := toInterfacesByStrings(previousResultRow)
		err := exFile.SetSheetRow(sheetName, axis, &sheetRow)
		if err != nil {
			logger.Errorf("Cannot set sheet row, got %v", err)
			return nil, err
		}
	}

	// 3. Set new result
	for rowId, row := range newResults {
		axis := fmt.Sprintf("A%v", rowId+dataIndexStart+len(previousFAPResults))
		sheetRow := row.ToInterfaces()
		err := exFile.SetSheetRow(sheetName, axis, &sheetRow)
		if err != nil {
			logger.Errorf("Cannot set sheet row, got %v", err)
			return nil, err
		}
	}

	// 4. Return
	return exFile.WriteToBuffer()
}

func makeResultFileByError(rows []dto.ErrorRow) (*bytes.Buffer, error) {
	const dataIndexStart = 3
	const numberOfColumn = 4
	exFile, sheetName := initFapResultSheetWithHeader()

	for rowId, row := range rows {
		rawDataRow := make([]interface{}, numberOfColumn)
		for index, rowCell := range row.RowData {
			if index <= numberOfColumn-2 {
				rawDataRow[index] = rowCell
			}
		}

		rawDataRow[numberOfColumn-1] = row.Reason
		axis := fmt.Sprintf("A%v", rowId+dataIndexStart)
		err := exFile.SetSheetRow(sheetName, axis, &rawDataRow)
		if err != nil {
			logger.Errorf("Cannot set sheet row, got %v", err)
			return nil, err
		}
	}

	return exFile.WriteToBuffer()
}

func initFapResultSheetWithHeader() (*excelize.File, string) {
	// 1. Init file
	exFile := excelize.NewFile()

	// 2. Create a new sheet & set active
	const sheetName = "Sheet1"
	index := exFile.NewSheet(sheetName)
	exFile.SetActiveSheet(index)

	// 3. Set header
	_ = exFile.SetSheetRow(sheetName, "A1", &[]interface{}{"Phone number (*)", "Points (*)", "Note", "Error"})
	_ = exFile.SetSheetRow(sheetName, "A2", &[]interface{}{"SĐT khách hàng", "Số điểm nạp", "Ghi chú giao dịch", "Kết quả"})

	// 4. return
	return exFile, sheetName
}

func toInterfacesByStrings(previousResultRow []string) []interface{} {
	var sheetRow []interface{}
	for i := 0; i < len(previousResultRow); i++ {
		sheetRow = append(sheetRow, previousResultRow[i])
	}
	return sheetRow
}
