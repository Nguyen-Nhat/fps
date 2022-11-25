package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
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
		field.String("row_data_raw"),
		field.Int32("task_index"),
		field.String("task_mapping"),
		field.String("task_depends_on"),
		field.String("task_request_raw"),
		field.String("task_response_raw"),
		field.Int16("status").Comment("Init=1; Processing=2; Failed=3; Finished=4"),
		field.Int16("error_display"),
	}
}
