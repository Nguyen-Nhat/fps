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
	ArrayItem    []*RequestFieldMD // optional, have value when type=array
	ArrayItemMap map[string]*RequestFieldMD

	// Constrains
	Required bool

	// Others
	ValueDependsOn       ValueDependsOn
	ValueDependsOnKey    string
	ValueDependsOnTaskID int
	Value                string // real value in string. RealValue = Value when Value is raw. RealValue is value after get from pattern in Value
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

func (ct ConfigTaskMD) Clone() ConfigTaskMD {
	requestParamsMap := make(map[string]*RequestFieldMD)
	for key, value := range ct.RequestParamsMap {
		reqField := value.Clone()
		requestParamsMap[key] = &reqField
	}

	requestBodyMap := make(map[string]*RequestFieldMD)
	for key, value := range ct.RequestBodyMap {
		reqField := value.Clone()
		requestBodyMap[key] = &reqField
	}

	return ConfigTaskMD{
		TaskIndex: ct.TaskIndex,
		// Request
		Endpoint:         ct.Endpoint,
		Method:           ct.Method,
		Header:           utils.CloneMap(ct.Header),
		RequestParamsMap: requestParamsMap,
		RequestBodyMap:   requestBodyMap,

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

func (rf RequestFieldMD) Clone() RequestFieldMD {
	return RequestFieldMD{
		Field:        rf.Field,
		Type:         rf.Type,
		ValuePattern: rf.ValuePattern,
		// Custom for array
		ArrayItem:    utils.CloneArray(rf.ArrayItem),
		ArrayItemMap: utils.CloneMap(rf.ArrayItemMap),
		// Constrains
		Required: rf.Required,
		// Others
		ValueDependsOn:       rf.ValueDependsOn,
		ValueDependsOnKey:    rf.ValueDependsOnKey,
		ValueDependsOnTaskID: rf.ValueDependsOnTaskID,
		Value:                rf.Value,
	}
}