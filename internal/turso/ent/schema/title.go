package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Title struct {
	ent.Schema
}

func (Title) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Title) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("slug").
			NotEmpty().
			Unique(),
		field.Enum("type").
			Values("FEATURE", "SERIES", "SPECIAL").
			Default("FEATURE"),
		field.Time("premiere_date"),
		field.Int("territories").
			Default(1).
			Positive(),
		field.String("current_master_version").
			Default(""),
		field.Enum("overall_status").
			Values("DRAFT", "ON_TRACK", "AT_RISK", "HOLD", "PROCESSING", "SHIPPED", "OVERDUE").
			Default("DRAFT"),
	}
}

func (Title) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("masters", Master.Type),
		edge.To("packages", MediaPackage.Type),
		edge.To("deliveries", Delivery.Type),
	}
}

func (Title) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("overall_status"),
		index.Fields("type"),
	}
}
