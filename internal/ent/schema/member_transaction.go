package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type MemberTransaction struct {
	ent.Schema
}

// Annotations of the User.
func (MemberTransaction) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "member_transaction"},
	}
}

// Fields of the FileAwardPoint.
func (MemberTransaction) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("file_award_point_id"),
		field.Int64("point"),
		field.String("phone").MaxLen(15).NotEmpty(),
		field.String("order_code").MaxLen(50).NotEmpty(),
		field.String("ref_id").MaxLen(50).NotEmpty(),
		field.Time("sent_time").Default(time.Now()),
		field.String("txn_desc").MaxLen(255),
		field.Int16("status").Default(0),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()),
	}
}
