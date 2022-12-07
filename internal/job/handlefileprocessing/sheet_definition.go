package handlefileprocessing

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils/excel/dto"
)

const (
	dataIndexStartInDataSheet = 2 // dữ liệu sheet data bắt đầu từ dòng số mấy

	sheetImportDataName = "Template"
	sheetMappingName    = "mapping"

	columnErrorName = "Lỗi"
)

var sheetMappingMetadata = dto.SheetMappingMetadata{
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
