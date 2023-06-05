package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ConfigTask struct {
	ent.Schema
}

// Annotations of the User.
func (ConfigTask) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "config_task"},
	}
}

// Fields of the FileAwardPoint.
func (ConfigTask) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("config_mapping_id").NonNegative(),
		field.Int32("task_index").Comment("For example: 1,2,3,..."),
		field.String("end_point").NotEmpty().Comment("For example: http://loyalty-core-api.loyalty-service/api/v1/grant"),
		field.String("method").NotEmpty().Comment("GET, POST, PUT, ..."),
		field.Text("header").Comment("Format JSON"),
		field.Text("request_params").NotEmpty().Comment("Format JSON"),
		field.Text("request_body").NotEmpty().Comment("Format JSON"),
		field.Int32("response_success_http_status").Comment("For example: http 200"),
		field.String("response_success_code_schema").Comment("Format JSON, contains path and values"),
		field.String("response_message_schema").Comment("Format JSON, contains path"),
		// default fields
		field.Time("created_at").Default(time.Now),
		field.String("created_by").NotEmpty(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
