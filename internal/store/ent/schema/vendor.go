package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Vendor holds the schema definition for the Vendor entity.
type Vendor struct {
	ent.Schema
}

// Fields of the Vendor.
func (Vendor) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.String("specialty").
			NotEmpty(), // e.g. "AUDIO_DUBBING", "SUBTITLES", "VIDEO_MASTERING"
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Vendor.
func (Vendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("packages", MediaPackage.Type),
	}
}
