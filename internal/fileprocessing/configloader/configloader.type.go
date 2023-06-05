package configloader

import "git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"

// ConfigMappingMD ...
type ConfigMappingMD struct {
	// metadata get from file processing
	DataStartAtRow     int      `json:"-"`
	RequireColumnIndex []string `json:"-"`
	ErrorColumnIndex   string   `json:"-"`

	// parameter in file
	FileParameters map[string]string `json:"-"`

	// List ConfigTaskMD
	Tasks []ConfigTaskMD
}

// ConfigTaskMD ...
type ConfigTaskMD struct {
	TaskIndex int
	// Request
	Endpoint         string
	Method           string
	Header           map[string]string
	RequestParamsMap map[string]*RequestFieldMD
	RequestBodyMap   map[string]*RequestFieldMD
	// Request that filled converted field's value
	RequestParams map[string]interface{}
	RequestBody   map[string]interface{}
	// Response
	Response ResponseMD
	// Row data in importing file -> is injected in validation phase
	ImportRowData  []string
	ImportRowIndex int
}

// RequestFieldMD ... Metadata for Request Field, use for describing RequestParams, RequestBody
type RequestFieldMD struct {
	Field        string // field name to request
	Type         string // support int, string, array (item is defined in array_item)
	ValuePattern string // it may be raw value or pattern to get

	// Custom for array
	ArrayItem *RequestFieldMDChild // optional, have value when type=array

	// Constrains
	Required bool

	// Others
	ValueDependsOn       ValueDependsOn
	ValueDependsOnKey    string
	ValueDependsOnTaskID int
	Value                string // real value in string. RealValue = Value when Value is raw. RealValue is value after get from pattern in Value
}

type RequestFieldMDChild struct {
	Field        string `json:"field"`        // field name to request
	Type         string `json:"type"`         // support int, string, array (item is defined in array_item)
	ValuePattern string `json:"valuePattern"` // it may be raw value or pattern to get

	// Others
	Value string `json:"value"` // real value in string. RealValue = Value when Value is raw. RealValue is value after get from pattern in Value
}

// ResponseMD ...
type ResponseMD struct {
	HttpStatusSuccess *int32
	Code              ResponseCode
	Message           ResponseMsg
}

// ---------------------------------------------------------------------------------------------------------------------

// ResponseCode ...
type ResponseCode struct {
	Path          string `json:"path"`
	SuccessValues string `json:"successValues"`
}

// ResponseMsg ...
type ResponseMsg struct {
	Path string `json:"path"`
}

func (ct *ConfigTaskMD) Clone() ConfigTaskMD {
	return ConfigTaskMD{
		TaskIndex: ct.TaskIndex,
		// Request
		Endpoint:         ct.Endpoint,
		Method:           ct.Method,
		Header:           utils.CloneMap(ct.Header),
		RequestParamsMap: utils.CloneMap(ct.RequestParamsMap),
		RequestBodyMap:   utils.CloneMap(ct.RequestBodyMap),
		// Request that filled converted field's value
		RequestParams: utils.CloneMap(ct.RequestParams),
		RequestBody:   utils.CloneMap(ct.RequestBody),
		// Response
		Response: ct.Response,
		// Row data in importing file -> is injected in validation phase
		ImportRowData:  ct.ImportRowData,
		ImportRowIndex: ct.ImportRowIndex,
	}
}
