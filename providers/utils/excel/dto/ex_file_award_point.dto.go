package dto

type FileAwardPointMetadata struct {
	Phone CellData[string]
	Point CellData[int]
	Note  CellData[string]
}

var _ Converter[FileAwardPointMetadata, FileAwardPointRow] = &FileAwardPointMetadata{}

// GetHeaders implement from Converter interface
func (f *FileAwardPointMetadata) GetHeaders() []string {
	var headers []string
	if len(f.Phone.ColumnName) > 0 {
		headers = append(headers, f.Phone.ColumnName)
	}
	if len(f.Point.ColumnName) > 0 {
		headers = append(headers, f.Point.ColumnName)
	}
	if len(f.Note.ColumnName) > 0 {
		headers = append(headers, f.Note.ColumnName)
	}
	return headers
}

// ToOutput implement from Converter interface
func (f *FileAwardPointMetadata) ToOutput(rowId int) FileAwardPointRow {
	return FileAwardPointRow{
		RowId: rowId,
		Phone: f.Phone.Value,
		Point: f.Point.Value,
		Note:  f.Note.Value,
	}
}

// GetMetadata implement from Converter interface
func (f *FileAwardPointMetadata) GetMetadata() *FileAwardPointMetadata {
	return f
}

type FileAwardPointRow struct {
	RowId int
	Phone string
	Point int
	Note  string
}
