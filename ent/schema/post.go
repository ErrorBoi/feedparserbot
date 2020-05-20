package schema

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

type TitleTranslations struct {
	RU string
	EN string
}

type SubjectTranslations struct {
	RU string
	EN string
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.JSON("title_translations", TitleTranslations{}),
		field.Text("subject").Nillable().Optional(),
		field.JSON("subject_translations", SubjectTranslations{}),
		field.String("url"),
		field.Time("published_at"),
		field.String("description"),
		field.String("h1"),
		field.Text("content"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Int("updated_by"),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("source", Source.Type).
			Ref("posts").
			Unique(),
	}
}
