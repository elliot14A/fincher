package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// WfResult records an evaluation, verdict, or decision produced during a workflow Run.
type WfResult struct {
	ent.Schema
}

func (WfResult) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (WfResult) Fields() []ent.Field {
	return []ent.Field{
		field.String("run_id").
			NotEmpty(),
		field.String("step_id").
			Optional(),
		field.String("judge").
			NotEmpty(),
		field.String("outcome").
			NotEmpty(),
		field.String("rationale").
			Default(""),
		field.Int("attempt").
			Default(1).
			Positive(),
	}
}

func (WfResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("run", Run.Type).
			Ref("results").
			Field("run_id").
			Required().
			Unique(),
		edge.From("step", Step.Type).
			Ref("results").
			Field("step_id").
			Unique(),
	}
}

func (WfResult) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("run_id"),
		index.Fields("step_id"),
		index.Fields("judge"),
		index.Fields("outcome"),
	}
}
