package dto

// Metadata ------------------------------------------------------------------------------------------------------------

type CellData[V any] struct {
	ColumnName string
	ValueRaw   string
	Value      V
	Constrains Constrains
}

type Constrains struct {
	IsRequired bool
	MinLength  *int
	MaxLength  *int
	Regexp     *string
}

//func (c *CellData[any]) format() {
//	// format for data ...
//}
//func (c *CellData[any]) validate() bool {
//	return false
//}

// Converter -----------------------------------------------------------------------------------------------------------

type Converter[META any, OUT any] interface {
	GetHeaders() []string
	GetMetadata() *META
	ToOutput(rowId int) OUT
}

// Out put -------------------------------------------------------------------------------------------------------------

type Sheet[R any] struct {
	DataIndexStart int
	Data           []R
	ErrorRows      []ErrorRow
}

type ErrorRow struct {
	RowId  int
	Reason string
}
