package chats

import (
	"context"

	"github.com/google/uuid"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	"github.com/elliot14A/fincher/internal/turso/ent/chatmessage"
	"github.com/elliot14A/fincher/internal/turso/ent/chatsession"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// CreateSession inserts a new, empty chat session.
func CreateSession(ctx context.Context, client *ent.Client, title string) domainerrors.Result[*models.ChatSession] {
	id := "chat-" + uuid.NewString()[:8]
	created, err := client.ChatSession.Create().
		SetID(id).
		SetTitle(title).
		Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.ChatSession](turso.MapEntError("chats.CreateSession", "chat_session", id, err))
	}
	return domainerrors.Ok(toDomainSession(created))
}

// GetSession fetches a single session with its messages ordered by sequence.
func GetSession(ctx context.Context, client *ent.Client, id string) domainerrors.Result[*models.ChatSession] {
	s, err := client.ChatSession.Query().
		Where(chatsession.IDEQ(id)).
		WithMessages(func(mq *ent.ChatMessageQuery) {
			mq.Order(ent.Asc(chatmessage.FieldSeq))
		}).
		Only(ctx)
	if err != nil {
		return domainerrors.Err[*models.ChatSession](turso.MapEntError("chats.GetSession", "chat_session", id, err))
	}
	return domainerrors.Ok(toDomainSession(s))
}

// ListSessions returns all sessions newest-first (without messages).
func ListSessions(ctx context.Context, client *ent.Client) domainerrors.Result[[]*models.ChatSession] {
	items, err := client.ChatSession.Query().
		Order(ent.Desc(chatsession.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return domainerrors.Err[[]*models.ChatSession](turso.MapEntError("chats.ListSessions", "chat_session", "", err))
	}
	return domainerrors.Ok(toDomainSessionList(items))
}

// UpdateTitle renames a session.
func UpdateTitle(ctx context.Context, client *ent.Client, id, title string) domainerrors.Result[*models.ChatSession] {
	updated, err := client.ChatSession.UpdateOneID(id).
		SetTitle(title).
		Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.ChatSession](turso.MapEntError("chats.UpdateTitle", "chat_session", id, err))
	}
	return domainerrors.Ok(toDomainSession(updated))
}

// DeleteSession removes a session and cascades to its messages.
func DeleteSession(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	if _, err := client.ChatMessage.Delete().
		Where(chatmessage.SessionIDEQ(id)).
		Exec(ctx); err != nil {
		return domainerrors.Err[bool](turso.MapEntError("chats.DeleteSession", "chat_message", id, err))
	}
	if err := client.ChatSession.DeleteOneID(id).Exec(ctx); err != nil {
		return domainerrors.Err[bool](turso.MapEntError("chats.DeleteSession", "chat_session", id, err))
	}
	return domainerrors.Ok(true)
}

// AppendMessage inserts a message into a session, assigning the next sequence number.
func AppendMessage(ctx context.Context, client *ent.Client, sessionID string, role models.ChatMessageRole, content string, citations []models.ChatCitation, sources []models.ChatSource) domainerrors.Result[*models.ChatMessage] {
	nextSeq, err := client.ChatMessage.Query().
		Where(chatmessage.SessionIDEQ(sessionID)).
		Count(ctx)
	if err != nil {
		return domainerrors.Err[*models.ChatMessage](turso.MapEntError("chats.AppendMessage", "chat_message", sessionID, err))
	}

	id := "msg-" + uuid.NewString()[:8]
	builder := client.ChatMessage.Create().
		SetID(id).
		SetSessionID(sessionID).
		SetRole(chatmessage.Role(role)).
		SetSeq(nextSeq + 1).
		SetContent(content)

	if meta := messageMetadata(citations, sources); meta != nil {
		builder.SetMetadata(meta)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		return domainerrors.Err[*models.ChatMessage](turso.MapEntError("chats.AppendMessage", "chat_message", id, err))
	}

	if err := client.ChatSession.UpdateOneID(sessionID).Exec(ctx); err != nil {
		return domainerrors.Err[*models.ChatMessage](turso.MapEntError("chats.AppendMessage", "chat_session", sessionID, err))
	}

	msg := toDomainMessage(created)
	return domainerrors.Ok(&msg)
}
