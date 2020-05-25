package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			Nillable().
			Optional(),
		field.Int("tg_id").
			Unique(),
		field.String("payment_info").
			Nillable().
			Optional(),
		field.Enum("role").
			Values("user", "editor", "admin").
			Default("user"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("settings", UserSettings.Type).
			Unique(),
		edge.To("sources", Source.Type),
	}
}
