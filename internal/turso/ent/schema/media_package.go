package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type MediaPackage struct {
	ent.Schema
}

func (MediaPackage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (MediaPackage) Fields() []ent.Field {
	return []ent.Field{
		field.String("title_id").
			NotEmpty(),
		field.Enum("component").
			Values("VIDEO", "AUDIO", "SUBTITLE", "METADATA"),
		field.String("language").
			NotEmpty(),
		field.String("version").
			NotEmpty(),
		field.String("vendor_id").
			NotEmpty(),
		field.String("derived_from_master_version").
			NotEmpty(),
		field.Int("redelivery_count").
			Default(0).
			NonNegative(),
		field.Enum("status").
			Values("PENDING", "VALID", "INVALIDATED", "RE_QC_PENDING").
			Default("PENDING"),
		field.String("market").
			Optional().
			Default(""),
	}
}

func (MediaPackage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("title", Title.Type).
			Ref("packages").
			Field("title_id").
			Required().
			Unique(),
		edge.From("vendor", Vendor.Type).
			Ref("packages").
			Field("vendor_id").
			Required().
			Unique(),
		edge.To("dependent_children", Dependency.Type),
		edge.To("parent_dependencies", Dependency.Type),
	}
}

func (MediaPackage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title_id"),
		index.Fields("vendor_id"),
		index.Fields("status"),
		index.Fields("component"),
		index.Fields("market"),
	}
}
