package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ChatSession is a persisted operator conversation with the read-only chat assistant.
type ChatSession struct {
	ent.Schema
}

func (ChatSession) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (ChatSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			Default(""),
	}
}

func (ChatSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", ChatMessage.Type),
	}
}

func (ChatSession) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
	}
}
