package chats

import (
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

const (
	citationsMetadataKey = "citations"
	sourcesMetadataKey   = "sources"
)

func messageMetadata(citations []models.ChatCitation, sources []models.ChatSource) map[string]any {
	meta := map[string]any{}
	if len(citations) > 0 {
		items := make([]any, len(citations))
		for i, c := range citations {
			items[i] = map[string]any{"label": c.Label, "source": c.Source, "tool": c.Tool}
		}
		meta[citationsMetadataKey] = items
	}
	if len(sources) > 0 {
		items := make([]any, len(sources))
		for i, s := range sources {
			items[i] = map[string]any{
				"kind":         s.Kind,
				"id":           s.ID,
				"display_name": s.DisplayName,
				"identifier":   s.Identifier,
				"subtitle":     s.Subtitle,
				"image_url":    s.ImageURL,
				"status":       s.Status,
			}
		}
		meta[sourcesMetadataKey] = items
	}
	if len(meta) == 0 {
		return nil
	}
	return meta
}

func citationsFromMetadata(metadata map[string]any) []models.ChatCitation {
	raw, ok := metadata[citationsMetadataKey]
	if !ok {
		return nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]models.ChatCitation, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, models.ChatCitation{
			Label:  asString(m["label"]),
			Source: asString(m["source"]),
			Tool:   asString(m["tool"]),
		})
	}
	return out
}

func sourcesFromMetadata(metadata map[string]any) []models.ChatSource {
	raw, ok := metadata[sourcesMetadataKey]
	if !ok {
		return nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]models.ChatSource, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, models.ChatSource{
			Kind:        asString(m["kind"]),
			ID:          asString(m["id"]),
			DisplayName: asString(m["display_name"]),
			Identifier:  asString(m["identifier"]),
			Subtitle:    asString(m["subtitle"]),
			ImageURL:    asString(m["image_url"]),
			Status:      asString(m["status"]),
		})
	}
	return out
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toDomainMessage(m *ent.ChatMessage) models.ChatMessage {
	if m == nil {
		return models.ChatMessage{}
	}
	return models.ChatMessage{
		Base: models.Base{
			ID:        m.ID,
			Metadata:  m.Metadata,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		SessionID: m.SessionID,
		Role:      models.ChatMessageRole(m.Role),
		Seq:       m.Seq,
		Content:   m.Content,
		Citations: citationsFromMetadata(m.Metadata),
		Sources:   sourcesFromMetadata(m.Metadata),
	}
}

func toDomainSession(s *ent.ChatSession) *models.ChatSession {
	if s == nil {
		return nil
	}
	session := &models.ChatSession{
		Base: models.Base{
			ID:        s.ID,
			Metadata:  s.Metadata,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		Title: s.Title,
	}
	if s.Edges.Messages != nil {
		session.Messages = make([]models.ChatMessage, len(s.Edges.Messages))
		for i, m := range s.Edges.Messages {
			session.Messages[i] = toDomainMessage(m)
		}
	}
	return session
}

func toDomainSessionList(items []*ent.ChatSession) []*models.ChatSession {
	out := make([]*models.ChatSession, len(items))
	for i, item := range items {
		out[i] = toDomainSession(item)
	}
	return out
}
