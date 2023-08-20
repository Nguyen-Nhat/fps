package fpRowGroup

// CreateRowGroupJob ...
type CreateRowGroupJob struct {
	FileID       int
	TaskIndex    int
	GroupByValue string
	TotalRows    int
	RowIndexList string
}
