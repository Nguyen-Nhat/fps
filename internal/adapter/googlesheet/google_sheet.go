package google_sheet

import (
	"context"
	"fmt"
	"net/http"

	"github.com/xuri/excelize/v2"
)

const googleSheetHost = "https://docs.google.com"

type GoogleSheetAdapter interface {
	GetXlsxFile(ctx context.Context, sheetId string) (*excelize.File, error)
}

type GoogleSheetImpl struct {
}

func NewGoogleSheetClient() GoogleSheetAdapter {
	return &GoogleSheetImpl{}
}

func (s *GoogleSheetImpl) GetXlsxFile(ctx context.Context, sheetId string) (*excelize.File, error) {
	exportUrl := fmt.Sprintf("%s/spreadsheets/d/%s/export", googleSheetHost, sheetId)
	// nolint
	resp, err := http.Get(exportUrl)
	if err != nil {
		return nil, err
	}
	return excelize.OpenReader(resp.Body)
}
