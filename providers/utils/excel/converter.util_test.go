package excel

import (
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"testing"
)

func Test_readDataReflect(t *testing.T) {
	metadata := dto.FileAwardPointMetadata{
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
	headerMap := map[string]int{"Phone number (*)": 0, "Points (*)": 1, "Note": 2}
	rowData := []string{"0393227489", "1000", "Note...."}
	b, c := readDataReflect(&metadata, headerMap, rowData)
	fmt.Printf("Res: %v - %v", b, c)
}
