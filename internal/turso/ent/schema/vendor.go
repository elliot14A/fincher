package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Vendor holds the schema definition for the Vendor entity.
type Vendor struct {
	ent.Schema
}

// Mixin of the Vendor.
func (Vendor) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Vendor.
func (Vendor) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("specialty").
			NotEmpty(),
	}
}

// Edges of the Vendor.
func (Vendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("packages", MediaPackage.Type),
	}
}
