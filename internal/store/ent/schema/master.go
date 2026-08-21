package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Master holds the schema definition for the Master entity.
type Master struct {
	ent.Schema
}

// Fields of the Master.
func (Master) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			NotEmpty(),
		field.String("title_id").
			NotEmpty(),
		field.String("version").
			NotEmpty(), // e.g. "V13"
		field.String("supersedes_version").
			Optional(), // e.g. "V12"
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Master.
func (Master) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("title", Title.Type).
			Ref("masters").
			Field("title_id").
			Required().
			Unique(),
	}
}
