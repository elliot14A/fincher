package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Title holds the schema definition for the Title entity.
type Title struct {
	ent.Schema
}

// Fields of the Title.
func (Title) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.Enum("type").
			Values("FEATURE", "SERIES", "SPECIAL").
			Default("FEATURE"),
		field.Time("premiere_date"),
		field.Int("territories").
			Default(1).
			Positive(),
		field.String("current_master_version").
			NotEmpty(),
		field.Enum("overall_status").
			Values("ON_TRACK", "AT_RISK", "HOLD", "PROCESSING", "SHIPPED").
			Default("PROCESSING"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Title.
func (Title) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("masters", Master.Type),
		edge.To("packages", MediaPackage.Type),
		edge.To("deliveries", Delivery.Type),
	}
}
