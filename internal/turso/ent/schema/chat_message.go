package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ChatMessage is a single turn (user or assistant) within a ChatSession.
// Assistant turns carry their tool-call citations inside the inherited metadata field.
type ChatMessage struct {
	ent.Schema
}

func (ChatMessage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (ChatMessage) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").
			NotEmpty(),
		field.Enum("role").
			Values("user", "assistant"),
		field.Int("seq").
			Default(1).
			Positive(),
		field.Text("content").
			Default(""),
	}
}

func (ChatMessage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", ChatSession.Type).
			Ref("messages").
			Field("session_id").
			Required().
			Unique(),
	}
}

func (ChatMessage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id"),
		index.Fields("session_id", "seq"),
	}
}
