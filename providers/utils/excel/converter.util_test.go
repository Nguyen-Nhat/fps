package excel

import (
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
	"testing"
)

func Test_readDataReflect(t *testing.T) {
	metadata := dto.SheetMappingMetadata{
		TaskId: dto.CellData[int]{
			ColumnName: "task_id",
			Constrains: dto.Constrains{IsRequired: true},
		},
		Endpoint: dto.CellData[string]{
			ColumnName: "endpoint",
			Constrains: dto.Constrains{IsRequired: true},
		},
		Header: dto.CellData[string]{
			ColumnName: "header",
			Constrains: dto.Constrains{IsRequired: false},
		},
		Request: dto.CellData[string]{
			ColumnName: "request",
			Constrains: dto.Constrains{IsRequired: true},
		},
		Response: dto.CellData[string]{
			ColumnName: "response",
			Constrains: dto.Constrains{IsRequired: true},
		},
	}

	headerMap := map[string]int{"task_id": 0, "endpoint": 1, "header": 2, "request": 3, "response": 4}
	rowData := []string{"1", "https://a.com", "abc", "egh"}
	b, c := readDataReflect(&metadata, headerMap, rowData)
	fmt.Printf("Res: %v - %v", b, c)
}
