package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Globalsettings holds the schema definition for the Globalsettings entity.
type Globalsettings struct {
	ent.Schema
}

// Fields of the Globalsettings.
func (Globalsettings) Fields() []ent.Field {
	return []ent.Field{
		field.Strings("clickbait_words").Optional(),
	}
}

// Edges of the Globalsettings.
func (Globalsettings) Edges() []ent.Edge {
	return nil
}
