package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Delivery struct {
	ent.Schema
}

func (Delivery) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Delivery) Fields() []ent.Field {
	return []ent.Field{
		field.String("title_id").
			NotEmpty(),
		field.String("country").
			NotEmpty(),
		field.Enum("status").
			Values("PENDING", "READY_TO_SHIP", "HOLD", "SHIPPED").
			Default("PENDING"),
		field.Time("target_date"),
	}
}

func (Delivery) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("title", Title.Type).
			Ref("deliveries").
			Field("title_id").
			Required().
			Unique(),
	}
}

func (Delivery) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title_id"),
		index.Fields("status"),
		index.Fields("country"),
	}
}
