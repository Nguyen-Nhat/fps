package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type FpsClient struct {
	ent.Schema
}

// Annotations of the User.
func (FpsClient) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "fps_client"},
	}
}

// Fields of the FileAwardPoint.
func (FpsClient) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("client_id"),
		field.String("name").NotEmpty(),
		field.String("description").NotEmpty(),
		field.String("sample_file_url").Default(""),
		field.String("import_file_template_url").Optional().Comment("URL of template file that client can download"),
		// default fields
		field.Time("created_at").Default(time.Now),
		field.String("created_by").NotEmpty(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
