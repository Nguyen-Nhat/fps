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
	ExcelMsgInvalidFile         = "file không hợp lệ"
)

const (
	EmptyString          = ""
	NilString            = "<nil>"
	SplitByDot           = "."
	SplitByDotChar       = '.'
	SplitByComma         = ","
	SplitByNewLine       = "\n"
	SplitByCommaAndSpace = ", "
	SplitByDeduce        = "->"
)

const (
	ParseInvalidSessionTokenCode = 209
)

const (
	ExtFileCSV     = "CSV"
	ExtFileXLSX    = "XLSX"
	ExtFileXLS     = "XLS"
	ExtFileUnknown = "unknown"
)

const (
	KafkaConsumerWithRetry               = "with-retry"
	KafkaConsumeTypeForUpdateResultAsync = "update-result-async"
)

const (
	Timeout = "timeout"
)

const (
	One = 1
	Two = 2

	MaxRetryDownload   = 3
	RetryDelayDownload = 5

	DefaultInsertBatchSize = 300
)

const (
	MessageSuccess = "success"
)

const (
	TypeString  = "string"
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeBoolean = "boolean"
	TypeJson    = "json"
)
