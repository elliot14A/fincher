package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.String("specialty").
			NotEmpty(),
	}
}

func (Vendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("packages", MediaPackage.Type),
	}
}

func (Vendor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("specialty"),
	}
}
