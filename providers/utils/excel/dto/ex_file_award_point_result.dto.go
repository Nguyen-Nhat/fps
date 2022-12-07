package dto

import "strconv"

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
func (f *FileAwardPointResultMetadata) ToOutput(rowId int) (FileAwardPointRow, error) {
	return FileAwardPointRow{
		RowId: rowId,
		Phone: f.Phone.Value,
		Point: f.Point.Value,
		Note:  f.Note.Value,
		Error: f.Error.Value,
	}, nil
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

func (f *FileAwardPointResultRow) ToInterfaces() []interface{} {
	var pointStr string
	pointStr = strconv.Itoa(f.Point)
	return []interface{}{f.Phone, pointStr, f.Note, f.Error}
}
