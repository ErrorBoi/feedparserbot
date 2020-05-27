package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// UserSettings holds the schema definition for the UserSettings entity.
type UserSettings struct {
	ent.Schema
}

// Fields of the UserSettings.
func (UserSettings) Fields() []ent.Field {
	return []ent.Field{
		field.Strings("urgent_words").Optional(),
		field.Strings("banned_words").Optional(),
		field.Enum("language").Values("RU", "EN").Default("RU"),
		field.Enum("sending_frequency").Values(
			"instant", "1h", "4h", "am", "pm",
			"mon", "tue", "wed", "thu", "fri", "sat", "sun").Default("instant"),
		field.Time("last_sending").Default(time.Now),
	}
}

// Edges of the UserSettings.
func (UserSettings) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("settings").
			Unique().
			Required(),
	}
}
