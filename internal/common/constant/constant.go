package constant

// GetListFileAwardPoint
const (
	PaginationMaxPage = 1000
	PaginationMaxSize = 200

	PaginationDefaultPage = 1
	PaginationDefaultSize = 10
)

// Excel error message
const (
	ExcelMsgInvalidFormat       = "sai định dạng"
	ExcelMsgMissOrInvalidFormat = "thiếu hoặc sai định dạng"
	ExcelMsgRequired            = "không được bỏ trống"
	ExcelMsgLength              = "đô dài phải"
	ExcelMsgValue               = "giá trị phải"
)

const (
	EmptyString          = ""
	SplitByDot           = "."
	SplitByNewLine       = "\n"
	SplitByCommaAndSpace = ", "
)

const (
	ParseInvalidSessionTokenCode = 209
)

const (
	ExtFileCSV     = "csv"
	ExtFileXLSX    = "xlsx"
	ExtFileUnknown = "unknown"
)

const (
	KafkaConsumerWithRetry               = "with-retry"
	KafkaConsumeTypeForUpdateResultAsync = "update-result-async"
)
