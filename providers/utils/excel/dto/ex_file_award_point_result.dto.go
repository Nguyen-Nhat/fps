package dto

type FileAwardPointResultMetadata struct {
	Phone CellData[string]
}

// GetHeaders implement from Converter interface
func (f *FileAwardPointResultMetadata) GetHeaders() []string {
	var headers []string
	if len(f.Phone.ColumnName) > 0 {
		headers = append(headers, f.Phone.ColumnName)
	}
	return headers
}

// ToOutput implement from Converter interface
func (f *FileAwardPointResultMetadata) ToOutput(rowId int) FileAwardPointResultRow {
	return FileAwardPointResultRow{
		RowId: rowId,
		Phone: f.Phone.Value,
		// ...
	}
}

// GetMetadata implement from Converter interface
func (f *FileAwardPointResultMetadata) GetMetadata() *FileAwardPointResultMetadata {
	return f
}

type FileAwardPointResultRow struct {
	RowId int
	Phone string
	Point int
	Note  string
	Error string
}
