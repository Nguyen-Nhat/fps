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
		field.Text("file_parameters").Comment("Format JSON. For storing parameters of client"),
		field.Int32("seller_id").Comment("seller id"),
		field.Int32("total_mapping").Default(0),
		field.Bool("need_group_row").Default(false),
		field.Int32("stats_total_row").Default(0),
		field.Int32("stats_total_processed").Default(0),
		field.Int32("stats_total_success").Default(0),
		field.String("error_display"),
		field.Time("created_at").Default(time.Now),
		field.String("created_by").NotEmpty(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
