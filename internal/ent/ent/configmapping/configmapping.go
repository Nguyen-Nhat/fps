// Code generated by ent, DO NOT EDIT.

package configmapping

import (
	"fmt"
	"time"
)

const (
	// Label holds the string label denoting the configmapping type in the database.
	Label = "config_mapping"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldClientID holds the string denoting the client_id field in the database.
	FieldClientID = "client_id"
	// FieldTotalTasks holds the string denoting the total_tasks field in the database.
	FieldTotalTasks = "total_tasks"
	// FieldDataStartAtRow holds the string denoting the data_start_at_row field in the database.
	FieldDataStartAtRow = "data_start_at_row"
	// FieldDataAtSheet holds the string denoting the data_at_sheet field in the database.
	FieldDataAtSheet = "data_at_sheet"
	// FieldRequireColumnIndex holds the string denoting the require_column_index field in the database.
	FieldRequireColumnIndex = "require_column_index"
	// FieldErrorColumnIndex holds the string denoting the error_column_index field in the database.
	FieldErrorColumnIndex = "error_column_index"
	// FieldResultFileConfig holds the string denoting the result_file_config field in the database.
	FieldResultFileConfig = "result_file_config"
	// FieldTimeout holds the string denoting the timeout field in the database.
	FieldTimeout = "timeout"
	// FieldInputFileType holds the string denoting the input_file_type field in the database.
	FieldInputFileType = "input_file_type"
	// FieldOutputFileType holds the string denoting the output_file_type field in the database.
	FieldOutputFileType = "output_file_type"
	// FieldMaxFileSize holds the string denoting the max_file_size field in the database.
	FieldMaxFileSize = "max_file_size"
	// FieldTenantID holds the string denoting the tenant_id field in the database.
	FieldTenantID = "tenant_id"
	// FieldUsingMerchantAttrName holds the string denoting the using_merchant_attr_name field in the database.
	FieldUsingMerchantAttrName = "using_merchant_attr_name"
	// FieldMerchantAttributeName holds the string denoting the merchant_attribute_name field in the database.
	FieldMerchantAttributeName = "merchant_attribute_name"
	// FieldUIConfig holds the string denoting the ui_config field in the database.
	FieldUIConfig = "ui_config"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldCreatedBy holds the string denoting the created_by field in the database.
	FieldCreatedBy = "created_by"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the configmapping in the database.
	Table = "config_mapping"
)

// Columns holds all SQL columns for configmapping fields.
var Columns = []string{
	FieldID,
	FieldClientID,
	FieldTotalTasks,
	FieldDataStartAtRow,
	FieldDataAtSheet,
	FieldRequireColumnIndex,
	FieldErrorColumnIndex,
	FieldResultFileConfig,
	FieldTimeout,
	FieldInputFileType,
	FieldOutputFileType,
	FieldMaxFileSize,
	FieldTenantID,
	FieldUsingMerchantAttrName,
	FieldMerchantAttributeName,
	FieldUIConfig,
	FieldCreatedAt,
	FieldCreatedBy,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// ClientIDValidator is a validator for the "client_id" field. It is called by the builders before save.
	ClientIDValidator func(int32) error
	// DefaultTotalTasks holds the default value on creation for the "total_tasks" field.
	DefaultTotalTasks int32
	// DefaultDataStartAtRow holds the default value on creation for the "data_start_at_row" field.
	DefaultDataStartAtRow int32
	// DataStartAtRowValidator is a validator for the "data_start_at_row" field. It is called by the builders before save.
	DataStartAtRowValidator func(int32) error
	// ResultFileConfigValidator is a validator for the "result_file_config" field. It is called by the builders before save.
	ResultFileConfigValidator func(string) error
	// DefaultTimeout holds the default value on creation for the "timeout" field.
	DefaultTimeout int32
	// DefaultInputFileType holds the default value on creation for the "input_file_type" field.
	DefaultInputFileType string
	// DefaultMaxFileSize holds the default value on creation for the "max_file_size" field.
	DefaultMaxFileSize int32
	// DefaultUsingMerchantAttrName holds the default value on creation for the "using_merchant_attr_name" field.
	DefaultUsingMerchantAttrName bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// CreatedByValidator is a validator for the "created_by" field. It is called by the builders before save.
	CreatedByValidator func(string) error
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)

// OutputFileType defines the type for the "output_file_type" enum field.
type OutputFileType string

// OutputFileType values.
const (
	OutputFileTypeXLSX OutputFileType = "XLSX"
	OutputFileTypeCSV  OutputFileType = "CSV"
)

func (oft OutputFileType) String() string {
	return string(oft)
}

// OutputFileTypeValidator is a validator for the "output_file_type" field enum values. It is called by the builders before save.
func OutputFileTypeValidator(oft OutputFileType) error {
	switch oft {
	case OutputFileTypeXLSX, OutputFileTypeCSV:
		return nil
	default:
		return fmt.Errorf("configmapping: invalid enum value for output_file_type field: %q", oft)
	}
}
