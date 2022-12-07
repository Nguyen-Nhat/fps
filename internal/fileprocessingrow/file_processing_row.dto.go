package fileprocessingrow

type CreateProcessingFileRowJob struct {
	FileId      int
	RowIndex    int
	RowDataRaw  string
	TaskIndex   int
	TaskMapping string
}

type UpdateAfterExecutingByJob struct {
	RequestRaw   string
	ResponseRaw  string
	Status       int16
	ErrorDisplay string
}
