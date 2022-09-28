package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
		field.Bool("active").
			Default(true),
		field.String("email"),
		field.String("phone_number"),
		field.String("password_hash"),
	}
}
