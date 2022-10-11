package excel

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
)

func LoadExcelByUrl(fileURL string) ([][]string, error) {
	data, err := loadDataFromUrl(fileURL)
	if err != nil {
		panic(err)
	}

	// Open the ZIP file with Excelize
	excel, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Reader", err)
		return nil, err
	}

	// Check no sheet
	lst := excel.GetSheetList()
	if len(lst) == 0 {
		fmt.Println("Empty document")
		return nil, fmt.Errorf("file empty")
	}

	// Get First Sheet
	sheetName := excel.GetSheetName(0)
	fmt.Printf("First sheet is %v\n", sheetName)
	sheetData, err := excel.GetRows(sheetName)

	// Close file excel
	defer func() {
		if err = excel.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Return
	return sheetData, nil
}

func loadDataFromUrl(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = r.Body.Close()
	}()

	return io.ReadAll(r.Body)
}
