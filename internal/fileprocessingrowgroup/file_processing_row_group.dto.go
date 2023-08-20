package fpRowGroup

// CreateRowGroupJob ...
type CreateRowGroupJob struct {
	FileID       int
	TaskIndex    int
	GroupByValue string
	TotalRows    int
	RowIndexList string
}

type UpdateAfterExecutingByJob struct {
	RequestCurl  string
	ResponseRaw  string
	Status       int16
	ErrorDisplay string
	ExecutedTime int64 // unit milliseconds
}
