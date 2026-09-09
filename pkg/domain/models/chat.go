package models

// ChatMessageRole enumerates the author of a chat message.
type ChatMessageRole string

const (
	ChatRoleUser      ChatMessageRole = "user"
	ChatRoleAssistant ChatMessageRole = "assistant"
)

// ChatCitation records a single tool call the assistant ran to ground an answer.
type ChatCitation struct {
	Label  string `json:"label"`
	Source string `json:"source"`
	Tool   string `json:"tool"`
}

// ChatSource is an enriched reference (e.g. a title with poster art) the answer discussed.
type ChatSource struct {
	Kind        string `json:"kind"`
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Identifier  string `json:"identifier"`
	Subtitle    string `json:"subtitle,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ChatMessage is a single persisted turn within a ChatSession.
type ChatMessage struct {
	Base
	SessionID string          `json:"session_id"`
	Role      ChatMessageRole `json:"role"`
	Seq       int             `json:"seq"`
	Content   string          `json:"content"`
	Citations []ChatCitation  `json:"citations,omitempty"`
	Sources   []ChatSource    `json:"sources,omitempty"`
}

// ChatSession is a persisted operator conversation with the chat assistant.
type ChatSession struct {
	Base
	Title    string        `json:"title"`
	Messages []ChatMessage `json:"messages,omitempty"`
}
