package xls

import (
	"bytes"
	"io"
	"net/http"

	"github.com/shakinm/xlsReader/xls"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/constant"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/common/errorz"
)

func LoadXlsByUrl(url, sheetName string) ([][]string, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = r.Body.Close()
	}()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	workbook, err := xls.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if sheetName == constant.EmptyString {
		sheet, err := workbook.GetSheet(0)
		if err != nil {
			return nil, err
		}
		return loadDataFromSheet(sheet), nil
	}

	for _, s := range workbook.GetSheets() {
		if s.GetName() == sheetName {
			return loadDataFromSheet(&s), nil
		}
	}

	return nil, errorz.SheetNotFound(sheetName)
}

func loadDataFromSheet(sheet *xls.Sheet) [][]string {
	var result [][]string
	rows := sheet.GetRows()
	for _, row := range rows {
		if row == nil {
			break
		}
		var dataRow []string
		isLastRow := true
		for _, cell := range row.GetCols() {
			dataRow = append(dataRow, cell.GetString())
			if cell.GetString() != constant.EmptyString {
				isLastRow = false
			}
		}
		if isLastRow {
			break
		}
		result = append(result, dataRow)
	}
	return result
}