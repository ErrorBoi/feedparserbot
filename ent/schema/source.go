package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Source holds the schema definition for the Source entity.
type Source struct {
	ent.Schema
}

// Fields of the Source.
func (Source) Fields() []ent.Field {
	return []ent.Field{
		field.String("url").Unique(),
		field.String("title"),
		field.Enum("language").Values("ru", "en").Default("ru"),
	}
}

// Edges of the Source.
func (Source) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Source.Type).
			From("parent").
			Unique(),
		edge.To("posts", Post.Type),
	}
}
