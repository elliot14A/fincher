package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Dependency struct {
	ent.Schema
}

func (Dependency) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			NotEmpty(),
		field.String("parent_id").
			NotEmpty(),
		field.String("child_id").
			NotEmpty(),
		field.Enum("dependency_type").
			Values("AUDIO_SYNC", "SUBTITLE_ALIGNMENT", "MASTER_DERIVATION"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

func (Dependency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("parent_package", MediaPackage.Type).
			Ref("dependent_children").
			Field("parent_id").
			Required().
			Unique(),
		edge.From("child_package", MediaPackage.Type).
			Ref("parent_dependencies").
			Field("child_id").
			Required().
			Unique(),
	}
}

func (Dependency) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_id", "child_id").
			Unique(),
		index.Fields("parent_id"),
		index.Fields("child_id"),
	}
}
