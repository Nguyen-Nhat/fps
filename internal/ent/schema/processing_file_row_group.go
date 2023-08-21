package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type ProcessingFileRowGroup struct {
	ent.Schema
}

// Annotations of the User.
func (ProcessingFileRowGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "processing_file_row_group"},
	}
}

// Fields of the FileAwardPoint.
func (ProcessingFileRowGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("file_id").Comment("File ID"),
		field.Int32("task_index").Comment("Task Index"),
		field.Text("group_by_value").Comment("Group by value, that is get from data in file excel"),
		field.Int32("total_rows").Comment("Total rows that have same group_by_value"),
		field.Text("row_index_list").Comment("List of row index, split by comma"),
		field.Text("group_request_curl").Comment("Request cURL"),
		field.Text("group_response_raw").Comment("Response raw"),
		field.Int16("status").Comment("Init=1; Processing=2; Failed=3; Success=4;"),
		field.String("error_display").Comment("Error for displaying"),
		field.Int64("executed_time").Comment("Executed time"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
