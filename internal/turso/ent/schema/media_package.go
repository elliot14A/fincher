package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MediaPackage holds the schema definition for the MediaPackage entity.
type MediaPackage struct {
	ent.Schema
}

// Mixin of the MediaPackage.
func (MediaPackage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the MediaPackage.
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
	}
}

// Edges of the MediaPackage.
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
