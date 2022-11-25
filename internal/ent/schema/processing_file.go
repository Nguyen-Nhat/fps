package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ProcessingFile struct {
	ent.Schema
}

// Annotations of the User.
func (ProcessingFile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "processing_file"},
	}
}

// Fields of the FileAwardPoint.
func (ProcessingFile) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("client_id"),
		field.String("display_name"),
		field.String("file_url"),
		field.String("result_file_url"),
		field.Int16("status").Comment("Init=1; Processing=2; Failed=3; Finished=4"),
		field.Int32("number_task_in_file"),
		field.Int32("stats_total_row"),
		field.Int32("stats_total_success"),
		field.String("created_by"),
	}
}
