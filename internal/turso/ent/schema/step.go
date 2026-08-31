package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Step tracks an individual node or phase execution in a Run.
type Step struct {
	ent.Schema
}

func (Step) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Step) Fields() []ent.Field {
	return []ent.Field{
		field.String("run_id").
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.Enum("status").
			Values("PENDING", "RUNNING", "COMPLETED", "FAILED", "SKIPPED").
			Default("RUNNING"),
		field.Time("started_at"),
		field.Time("ended_at").
			Optional().
			Nillable(),
	}
}

func (Step) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("run", Run.Type).
			Ref("steps").
			Field("run_id").
			Required().
			Unique(),
		edge.To("results", WfResult.Type),
	}
}

func (Step) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("run_id"),
		index.Fields("status"),
		index.Fields("name"),
	}
}
