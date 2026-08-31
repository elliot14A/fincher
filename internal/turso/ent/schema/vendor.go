package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Vendor struct {
	ent.Schema
}

func (Vendor) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Vendor) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.JSON("components", []string{}).
			Comment("Covered media components (VIDEO, AUDIO, SUBTITLE, METADATA)"),
		field.JSON("markets", []string{}).
			Comment("Covered language-market tags (e.g. en-US, de-DE, hi-IN, te-IN)"),
		field.Float("hourly_rate_usd").
			Default(0.0).
			Min(0.0),
		field.Int("turnaround_hours").
			Default(24).
			Positive(),
	}
}

func (Vendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("packages", MediaPackage.Type),
	}
}

func (Vendor) Indexes() []ent.Index {
	return []ent.Index{}
}
