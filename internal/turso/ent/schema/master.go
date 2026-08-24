package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Master struct {
	ent.Schema
}

func (Master) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			NotEmpty(),
		field.String("title_id").
			NotEmpty(),
		field.String("version").
			NotEmpty(),
		field.String("supersedes_version").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

func (Master) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("title", Title.Type).
			Ref("masters").
			Field("title_id").
			Required().
			Unique(),
	}
}

func (Master) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title_id"),
		index.Fields("title_id", "version").Unique(),
	}
}
