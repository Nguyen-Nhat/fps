package dto

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/membertxn"
)

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
func (f *FileAwardPointMetadata) ToOutput(rowId int) (FileAwardPointRow, error) {
	return FileAwardPointRow{
		RowId: rowId,
		Phone: f.Phone.Value,
		Point: f.Point.Value,
		Note:  f.Note.Value,
	}, nil
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
	Error string
}

func (f FileAwardPointRow) SetError(errMessage string) {

}

func MapMemberTxnToFileAwardPointRow(records []membertxn.MemberTransaction) []membertxn.MemberTxnDTO {
	var rs []membertxn.MemberTxnDTO
	for _, record := range records {
		rs = append(rs, membertxn.MemberTxnDTO{
			ID:               int64(record.ID),
			FileAwardPointID: int64(record.FileAwardPointID),
			Point:            record.Point,
			Phone:            record.Phone,
			RefID:            record.RefID,
			TxnDesc:          record.TxnDesc,
			Status:           record.Status,
			Error:            record.Error,
			LoyaltyTxnID:     record.LoyaltyTxnID,
		})
	}
	return rs
}
