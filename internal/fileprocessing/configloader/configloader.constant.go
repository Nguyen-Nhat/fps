package configloader

const (
	TypeString = "STRING"
	TypeInt    = "INT"
	TypeLong   = "LONG"
	TypeDouble = "DOUBLE"
	TypeArray  = "ARRAY"
)

type ValueDependsOn string

const (
	ValueDependsOnNone  ValueDependsOn = "NONE"  // No depend
	ValueDependsOnExcel ValueDependsOn = "EXCEL" // Excel column in importing file
	ValueDependsOnParam ValueDependsOn = "PARAM" // Parameters that is give when submit file
	ValueDependsOnTask  ValueDependsOn = "TASK"  // Depends on response of previous task (request)
)

const (
	prefixMappingRequest          = "$"
	prefixMappingRequestParameter = "$param"
	prefixMappingRequestResponse  = "$response"
)
