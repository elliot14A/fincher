package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Delivery holds the schema definition for the Delivery entity.
type Delivery struct {
	ent.Schema
}

// Mixin of the Delivery.
func (Delivery) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Delivery.
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

// Edges of the Delivery.
func (Delivery) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("title", Title.Type).
			Ref("deliveries").
			Field("title_id").
			Required().
			Unique(),
	}
}
