package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Run tracks an agent workflow execution.
type Run struct {
	ent.Schema
}

func (Run) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Run) Fields() []ent.Field {
	return []ent.Field{
		field.String("title_slug").
			Default("GLOBAL"),
		field.String("trigger").
			NotEmpty(),
		field.Enum("status").
			Values("PENDING", "RUNNING", "COMPLETED", "FAILED", "ESCALATED").
			Default("RUNNING"),
		field.Time("started_at"),
		field.Time("ended_at").
			Optional().
			Nillable(),
	}
}

func (Run) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("steps", Step.Type),
		edge.To("results", WfResult.Type),
	}
}

func (Run) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("trigger"),
		index.Fields("title_slug"),
	}
}
