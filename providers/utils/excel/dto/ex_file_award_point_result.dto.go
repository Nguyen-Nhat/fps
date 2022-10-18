package dto

type FileAwardPointResultMetadata struct {
	Phone CellData[string]
	Point CellData[int]
	Note  CellData[string]
	Error CellData[string]
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
func (f *FileAwardPointResultMetadata) ToOutput(rowId int) FileAwardPointRow {
	return FileAwardPointRow{
		RowId: rowId,
		Phone: f.Phone.Value,
		Point: f.Point.Value,
		Note:  f.Note.Value,
		Error: f.Error.Value,
	}
}

// GetMetadata implement from Converter interface
func (f *FileAwardPointResultMetadata) GetMetadata() *FileAwardPointResultMetadata {
	return f
}
