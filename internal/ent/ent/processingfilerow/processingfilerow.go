// Code generated by ent, DO NOT EDIT.

package processingfilerow

const (
	// Label holds the string label denoting the processingfilerow type in the database.
	Label = "processing_file_row"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFileID holds the string denoting the file_id field in the database.
	FieldFileID = "file_id"
	// FieldRowIndex holds the string denoting the row_index field in the database.
	FieldRowIndex = "row_index"
	// FieldRowDataRaw holds the string denoting the row_data_raw field in the database.
	FieldRowDataRaw = "row_data_raw"
	// FieldTaskIndex holds the string denoting the task_index field in the database.
	FieldTaskIndex = "task_index"
	// FieldTaskMapping holds the string denoting the task_mapping field in the database.
	FieldTaskMapping = "task_mapping"
	// FieldTaskDependsOn holds the string denoting the task_depends_on field in the database.
	FieldTaskDependsOn = "task_depends_on"
	// FieldTaskRequestRaw holds the string denoting the task_request_raw field in the database.
	FieldTaskRequestRaw = "task_request_raw"
	// FieldTaskResponseRaw holds the string denoting the task_response_raw field in the database.
	FieldTaskResponseRaw = "task_response_raw"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldErrorDisplay holds the string denoting the error_display field in the database.
	FieldErrorDisplay = "error_display"
	// Table holds the table name of the processingfilerow in the database.
	Table = "processing_file_row"
)

// Columns holds all SQL columns for processingfilerow fields.
var Columns = []string{
	FieldID,
	FieldFileID,
	FieldRowIndex,
	FieldRowDataRaw,
	FieldTaskIndex,
	FieldTaskMapping,
	FieldTaskDependsOn,
	FieldTaskRequestRaw,
	FieldTaskResponseRaw,
	FieldStatus,
	FieldErrorDisplay,
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