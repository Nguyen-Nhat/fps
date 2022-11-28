package schema

import (
	"time"

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
		field.Int32("client_id"),
		field.String("display_name").NotEmpty(),
		field.String("file_url").NotEmpty(),
		field.String("result_file_url"),
		field.Int16("status").Comment("Init=1; Processing=2; Failed=3; Finished=4"),
		field.Int32("total_mapping"),
		field.Int32("stats_total_row"),
		field.Int32("stats_total_success"),
		field.Time("created_at").Default(time.Now),
		field.String("created_by").NotEmpty(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
