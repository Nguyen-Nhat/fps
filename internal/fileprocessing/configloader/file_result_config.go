package configloader

// ResultFileConfigMD ... this metadata describes the value will be added into a ColumnKey in Result File
// FieldPath helps us to get value from API 's response
type ResultFileConfigMD struct {
	ColumnKey     string `json:"column_key"` // E.g: A, B, AB, ...
	ValuePath     string `json:"value_path"`
	ValueInTaskID int32  `json:"value_in_task_id"`
}
