package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type FileAwardPoint struct {
	ent.Schema
}

// Annotations of the User.
func (FileAwardPoint) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "file_award_point"},
	}
}

// Fields of the FileAwardPoint.
func (FileAwardPoint) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("merchant_id"),
		field.String("display_name").NotEmpty(),
		field.String("file_url").NotEmpty(),
		field.String("result_file_url").NotEmpty(),
		field.Int16("status").Default(0),
		field.Int16("stats_total_row").Default(0),
		field.Int16("stats_total_success").Default(0),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()).UpdateDefault(time.Now()),
		field.String("created_by"),
		field.String("updated_by"),
	}
}
