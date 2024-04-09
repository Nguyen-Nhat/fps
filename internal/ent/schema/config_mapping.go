package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type ConfigMapping struct {
	ent.Schema
}

// Annotations of the User.
func (ConfigMapping) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "config_mapping"},
	}
}

// Fields of the FileAwardPoint.
func (ConfigMapping) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("client_id").NonNegative(),
		field.Int32("total_tasks").Default(0),
		field.Int32("data_start_at_row").
			Default(0).NonNegative().Min(0).Comment("Data start in this row index. Apply for excel, csv"),
		field.String("data_at_sheet").Comment("Default is first sheet in file"),
		field.String("require_column_index").Comment("For example: A,B,C"),
		field.String("error_column_index").Comment("Index of column, that FPS will fill when error happens"),
		field.Int32("timeout").Default(86400).Comment("Time out of template in seconds (default 24h as 86400 seconds)"),
		field.String("input_file_type").Default("XLSX").Comment("Các định dạng cho phép của file input, cách nhau bằng dấu phẩy (ex: \"XLSX,CSV\")"),
		field.Enum("output_file_type").Values("XLSX", "CSV").Default("XLSX").Comment("Type of file output"),
		// default fields
		field.Time("created_at").Default(time.Now),
		field.String("created_by").NotEmpty(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
