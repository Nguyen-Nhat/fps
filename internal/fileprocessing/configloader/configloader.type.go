package configloader

import (
	"encoding/json"

	customFunc "git.teko.vn/loyalty-system/loyalty-file-processing/pkg/customfunction/common"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
)

// ConfigMappingMD ...
type ConfigMappingMD struct {
	// metadata get from file processing
	DataStartAtRow     int      `json:"-"`
	DataAtSheet        string   `json:"-"`
	RequireColumnIndex []string `json:"-"`
	ErrorColumnIndex   string   `json:"-"`

	// parameter in file
	FileParameters map[string]interface{} `json:"-"`

	// List ConfigTaskMD
	Tasks []ConfigTaskMD
	//Task *ConfigTaskMD // todo use Task instead of Tasks, because Tasks should be used as a metadata, and Task will contains data of row
	RowGroupValue string `json:"-"` // todo remove when use Task, RowGroupValue only support group value for 1 task, not multi-task, need cause this
}

// GetConfigTaskMD ...  always return task
func (cf *ConfigMappingMD) GetConfigTaskMD(taskIndex int) ConfigTaskMD {
	for _, t := range cf.Tasks {
		if t.TaskIndex == taskIndex {
			return t
		}
	}

	return ConfigTaskMD{}
}

func (cf *ConfigMappingMD) IsSupportGrouping() bool {
	for _, task := range cf.Tasks {
		if task.RowGroup.IsSupportGrouping() { // at least one task support grouping
			return true
		}
	}
	return false
}

// ConfigTaskMD ...
type ConfigTaskMD struct {
	TaskIndex int
	TaskName  string
	// Request
	Endpoint         string
	Method           string
	RequestHeaderMap map[string]*RequestFieldMD
	PathParamsMap    map[string]*RequestFieldMD
	RequestParamsMap map[string]*RequestFieldMD
	RequestBodyMap   map[string]*RequestFieldMD
	// Request that filled converted field's value
	RequestHeader map[string]interface{}
	PathParams    map[string]interface{}
	RequestParams map[string]interface{}
	RequestBody   map[string]interface{}
	// Response
	Response ResponseMD
	// Group
	RowGroup RowGroupMD
	// Row data in importing file -> is injected in validation phase
	ImportRowHeader []string
	ImportRowData   []string
	ImportRowIndex  int
}

// RequestFieldMD ... Metadata for Request Field, use for describing RequestParams, RequestBody
type RequestFieldMD struct {
	Field        string // field name to request
	Type         string // support int, string, array (item is defined in array_item)
	ValuePattern string // it may be raw value or pattern to get

	// For array
	ArrayItem    []*RequestFieldMD // optional, have value when type=array
	ArrayItemMap map[string]*RequestFieldMD

	// For nested object
	Items    []*RequestFieldMD // optional, have value when type=object
	ItemsMap map[string]*RequestFieldMD

	// Constrains
	Required bool

	// Others
	ValueDependsOn       ValueDependsOn
	ValueDependsOnKey    string
	ValueDependsOnTaskID int
	ValueDependsOnFunc   customFunc.CustomFunction
	Value                string // real value in string. RealValue = Value when Value is raw. RealValue is value after get from pattern in Value
}

// ResponseMD ...
type ResponseMD struct {
	HttpStatusSuccess *int32
	Code              ResponseCode
	Message           ResponseMsg
	MessageTransforms map[int]MessageTransformation
}

// ---------------------------------------------------------------------------------------------------------------------

// ResponseCode ...
type ResponseCode struct {
	Path          string `json:"path"`
	SuccessValues string `json:"successValues"`
	// MustHaveValueInPath ... this field is temporary, we will define more general rule later
	MustHaveValueInPath string `json:"mustHaveValueInPath"`
}

// ResponseMsg ...
type ResponseMsg struct {
	Path string `json:"path"`
}

type MessageTransformation struct {
	HttpStatus int    `json:"httpStatus"`
	Message    string `json:"message"`
	// We can support multi-language by using `messageTranslated`
	// E.g: { "vi": "Đây là message", "en": "This is a message" }
	MessageTranslated any `json:"messageTranslated"`
}

// RowGroupMD ...
type RowGroupMD struct {
	GroupByColumnsRaw string `json:"-"`
	GroupByColumns    []int  `json:"groupByColumns"`
	GroupSizeLimit    int    `json:"-"`
}

func (rg RowGroupMD) IsSupportGrouping() bool {
	return len(rg.GroupByColumns) > 0
}

func (ct ConfigTaskMD) Clone() ConfigTaskMD {
	headerMap := make(map[string]*RequestFieldMD)
	for key, value := range ct.RequestHeaderMap {
		reqField := value.Clone()
		headerMap[key] = &reqField
	}

	pathParamsMap := make(map[string]*RequestFieldMD)
	for key, value := range ct.PathParamsMap {
		reqField := value.Clone()
		pathParamsMap[key] = &reqField
	}

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
		TaskName:  ct.TaskName,
		// Request
		Endpoint:         ct.Endpoint,
		Method:           ct.Method,
		RequestHeader:    utils.CloneMap(ct.RequestHeader),
		RequestHeaderMap: headerMap,
		PathParamsMap:    pathParamsMap,
		RequestParamsMap: requestParamsMap,
		RequestBodyMap:   requestBodyMap,

		// Request that filled converted field's value
		PathParams:    utils.CloneMap(ct.PathParams),
		RequestParams: utils.CloneMap(ct.RequestParams),
		RequestBody:   utils.CloneMap(ct.RequestBody),

		// Response
		Response: ct.Response,

		// Row Group
		RowGroup: ct.RowGroup,

		// Row data in importing file -> is injected in validation phase
		ImportRowHeader: ct.ImportRowHeader,
		ImportRowData:   ct.ImportRowData,
		ImportRowIndex:  ct.ImportRowIndex,
	}
}

func (rf RequestFieldMD) Clone() RequestFieldMD {
	js, _ := json.Marshal(rf)
	res := RequestFieldMD{}
	_ = json.Unmarshal(js, &res)
	return res
}
