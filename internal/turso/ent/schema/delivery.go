package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Delivery holds the schema definition for the Delivery entity.
type Delivery struct {
	ent.Schema
}

// Fields of the Delivery.
func (Delivery) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			NotEmpty(),
		field.String("title_id").
			NotEmpty(),
		field.String("country").
			NotEmpty(), // ISO-3166-1 alpha-2 or country code, e.g. "US", "ES", "JP"
		field.Enum("status").
			Values("PENDING", "READY_TO_SHIP", "HOLD", "SHIPPED").
			Default("PENDING"),
		field.Time("target_date"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
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
