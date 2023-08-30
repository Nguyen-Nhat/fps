package fileprocessingrow

type CreateProcessingFileRowJob struct {
	FileID       int
	RowIndex     int
	RowDataRaw   string
	TaskIndex    int
	TaskMapping  string
	GroupByValue string
}

type UpdateAfterExecutingByJob struct {
	Task ProcessingFileRow // original data

	TaskMapping  string
	RequestCurl  string
	ResponseRaw  string
	Status       int16
	ErrorDisplay string
	ExecutedTime int64 // unit milliseconds
}

type StatisticData struct {
	IsFinished    bool
	ErrorDisplays map[int]string

	TotalProcessed int
	TotalSuccess   int
	TotalFailed    int
	TotalWaiting   int
}
