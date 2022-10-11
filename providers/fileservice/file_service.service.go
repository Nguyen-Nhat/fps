package fileservice

import (
	"bytes"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"github.com/xuri/excelize/v2"
	"strconv"
)

type (
	// IService ...
	IService interface {
		UploadFileAwardPointResult([]dto.FileAwardPointResultRow, string) (string, error)
		LoadFileAwardPoint(string) (*dto.Sheet[dto.FileAwardPointRow], error)
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

func (s *Service) UploadFileAwardPointResult(fabResults []dto.FileAwardPointResultRow, previousResultFileUrl string) (string, error) {
	// 1. Load previous result file
	logger.Infof("fetch file %v", previousResultFileUrl)
	var previousFAPResults []dto.FileAwardPointResultRow

	// 2. Mix previous result with current result
	// todo need to check this case, we can append data (as below line), override data or mix data
	combinedFAPResults := append(previousFAPResults, fabResults...)

	// 3. Convert data to bytes
	dataByteBuffer, err := makeResultFile(combinedFAPResults)
	if err != nil {
		logger.Errorf("Failed to convert data to excel %v", err)
	}

	// 4. Build request
	fileName := getResultFileNameFromUrl(previousResultFileUrl)
	req := uploadFileRequest{
		FileData: dataByteBuffer.Bytes(),
		FileName: fileName,
	}

	// 5. Request upload
	res, err := s.client.uploadFile(req)
	if err != nil {
		logger.Errorf("Upload file to File Service failed: %v", err)
		return "", err
	}

	// 6. Return
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

func getResultFileNameFromUrl(previousResultFileUrl string) string {
	fileNameMatch := constant.FileNameRegex.FindStringSubmatch(previousResultFileUrl)
	var fileName string
	if len(fileNameMatch) < 2 {
		fileName = "file_award_point_result.xlsx"
	}
	fileName = fileNameMatch[1]
	return fileName
}

func makeResultFile(rows []dto.FileAwardPointResultRow) (*bytes.Buffer, error) {
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
		_ = f.SetSheetRow(sheetName, axis, &[]interface{}{row.RowId, row.Phone, pointStr, row.Note, row.Error})
	}

	return f.WriteToBuffer()
}
