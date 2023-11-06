package configloader

const (
	TypeString  = "string"
	TypeInteger = "integer" // include int32 and int64
	TypeNumber  = "number"  // include float32 and float64
	TypeBoolean = "boolean"

	TypeArray  = "array"
	TypeObject = "object"

	TypeJson = "json"
)

type ValueDependsOn string

const (
	ValueDependsOnNone  ValueDependsOn = "NONE"  // No depend
	ValueDependsOnExcel ValueDependsOn = "EXCEL" // Excel column in importing file
	ValueDependsOnParam ValueDependsOn = "PARAM" // Parameters that is give when submit file
	ValueDependsOnTask  ValueDependsOn = "TASK"  // Depends on response of previous task (request)
	ValueDependsOnFunc  ValueDependsOn = "FUNC"  // Depends on function, and parameters of this function can be depended on Excel/Param
)

const (
	PrefixMappingRequest          = "$"
	PrefixMappingRequestParameter = "$param"
	prefixMappingRequestResponse  = "$response"
)
