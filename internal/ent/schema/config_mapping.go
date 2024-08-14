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
		field.String("result_file_config").Optional().MaxLen(500).Comment("JSON string, config new column in file result for display process result"),
		field.Int32("timeout").Default(86400).Comment("Time out of template in seconds (default 24h as 86400 seconds)"),
		field.String("input_file_type").Default("XLSX").Comment("Các định dạng cho phép của file input, cách nhau bằng dấu phẩy (ex: \"XLS,XLSX,CSV\")"),
		field.Enum("output_file_type").Values("XLSX", "CSV").Optional().Comment("Type of file output (XLS, XLSX, CSV). If null, output type is input type. If has value will force output type"),
		field.Int32("max_file_size").Default(5).Comment("Max file size (MB) that client can upload"),
		field.String("tenant_id").Optional().Comment("Tenant Id of client"),
		field.Bool("using_merchant_attr_name").Default(false).Comment("If 1, when import/get data, FPS will filter by sellerId, platformId,... (based on merchant_attribute_name value)"),
		field.String("merchant_attribute_name").Optional().Comment("Attribute name of users attribute that is used for filtering data"),
		field.Text("ui_config").Optional().Comment("UI config for client. Eg: show hide elements, change positions, ... Ref: https://confluence.teko.vn/display/SupplyChain/%5BFPS%5D+UI+Config"),
		// default fields
		field.Time("created_at").Default(time.Now),
		field.String("created_by").NotEmpty(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
