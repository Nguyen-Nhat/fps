package excel

import (
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"testing"
)

func Test_openFileURL_convertToStruct(t *testing.T) {
	fullURLFile := "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/9/Fumart%20Loyalty%20-%20nap%20diem%20KH.xlsx"

	fileAwardPointMetadata := dto.FileAwardPointMetadata{
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

	sheetData, err := LoadExcelByUrl(fullURLFile)
	if err != nil {
		return
	}

	dataIndexStart := 3
	sheet, err := ConvertToStruct[
		dto.FileAwardPointMetadata,
		dto.FileAwardPointRow,
		dto.Converter[dto.FileAwardPointMetadata, dto.FileAwardPointRow],
	](dataIndexStart, &fileAwardPointMetadata, sheetData)
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	fmt.Printf("done %v", sheet)
}

func Test_openFileURL_convertToStruct2(t *testing.T) {
	fullURLFile := "https://storage.googleapis.com/develop-teko-storage/media/doc/2022/10/9/Fumart%20Loyalty%20-%20nap%20diem%20KH.xlsx"

	phoneMinLength := 9
	phoneMaxLength := 9
	phoneRegexp := "^[0-9]+$" // "^[a-z]+$"
	fileAwardPointMetadata := dto.FileAwardPointResultMetadata{
		Phone: dto.CellData[string]{
			ColumnName: "Phone number (*)",
			Constrains: dto.Constrains{
				IsRequired: true,
				MinLength:  &phoneMinLength,
				MaxLength:  &phoneMaxLength,
				Regexp:     &phoneRegexp,
			},
		},
	}

	sheetData, err := LoadExcelByUrl(fullURLFile)
	if err != nil {
		return
	}

	dataIndexStart := 3
	sheet, err := ConvertToStruct[
		dto.FileAwardPointResultMetadata,
		dto.FileAwardPointResultRow,
		dto.Converter[dto.FileAwardPointResultMetadata, dto.FileAwardPointResultRow],
	](dataIndexStart, &fileAwardPointMetadata, sheetData)
	if err != nil {
		fmt.Printf("error %v", err)
		return
	}
	fmt.Printf("done %v", sheet)
}
