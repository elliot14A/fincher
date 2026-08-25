package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Upload holds the schema definition for the Upload entity.
type Upload struct {
	ent.Schema
}

// Fields of the Upload.
func (Upload) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable(),
		field.String("filename").
			NotEmpty(),
		field.String("mime_type").
			NotEmpty(),
		field.Bytes("data").
			NotEmpty(),
		field.Int64("size_bytes").
			Positive(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Upload.
func (Upload) Edges() []ent.Edge {
	return nil
}
