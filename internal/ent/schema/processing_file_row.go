package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type ProcessingFileRow struct {
	ent.Schema
}

// Annotations of the User.
func (ProcessingFileRow) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "processing_file_row"},
	}
}

// Fields of the FileAwardPoint.
func (ProcessingFileRow) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("file_id"),
		field.Int32("row_index"),
		field.Text("row_data_raw"),
		field.Int32("task_index"),
		field.Text("task_mapping"),
		field.String("task_depends_on"),
		field.Text("group_by_value").Default(""),
		field.Text("task_request_curl"),
		field.Text("task_request_raw"),
		field.Text("task_response_raw"),
		field.Int16("status").Comment("Init=1; ; Failed=3; Success=4; Timeout=5"),
		field.String("error_display"),
		field.Int64("executed_time"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
