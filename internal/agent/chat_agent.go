package agent

import (
	"context"
	"database/sql"
	"strings"

	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
	"github.com/elliot14A/fincher/pkg/mcp"
	"github.com/elliot14A/fincher/prompts"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
)

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
)

type ChatTurn struct {
	Role    ChatRole `json:"role"`
	Content string   `json:"content"`
}

type ChatCitation struct {
	Label  string `json:"label"`
	Source string `json:"source"`
	Tool   string `json:"tool"`
}

type ChatSource struct {
	Kind        string `json:"kind"`
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Identifier  string `json:"identifier"`
	Subtitle    string `json:"subtitle,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	Status      string `json:"status,omitempty"`
}

type ChatAnswer struct {
	Answer    string         `json:"answer"`
	Citations []ChatCitation `json:"citations"`
	Sources   []ChatSource   `json:"sources"`
}

func buildChatTools(tursoDB *sql.DB, mcpClient *mcp.Client) ([]tool.Tool, error) {
	tursoTool, err := tools.NewTursoQueryTool(tursoDB)
	if err != nil {
		return nil, err
	}
	clickhouseTool, err := tools.NewClickhouseQueryTool(mcpClient)
	if err != nil {
		return nil, err
	}
	return []tool.Tool{tursoTool, clickhouseTool}, nil
}

func renderChatHistory(history []ChatTurn, userMessage string) string {
	var b strings.Builder
	for _, t := range history {
		label := "Operator"
		if t.Role == ChatRoleAssistant {
			label = "Assistant"
		}
		b.WriteString(label)
		b.WriteString(": ")
		b.WriteString(strings.TrimSpace(t.Content))
		b.WriteString("\n\n")
	}
	b.WriteString("Operator: ")
	b.WriteString(strings.TrimSpace(userMessage))
	return b.String()
}

func citationSource(toolName string) string {
	if toolName == "query_clickhouse" {
		return "clickhouse"
	}
	return "sqlite"
}

func toolCallLabel(fc *genai.FunctionCall) string {
	if fc.Args != nil {
		if raw, ok := fc.Args["sql"]; ok {
			if s, ok := raw.(string); ok && strings.TrimSpace(s) != "" {
				return strings.TrimSpace(s)
			}
		}
	}
	return fc.Name
}

func AnswerQuery(
	ctx context.Context,
	m model.LLM,
	tursoClient *ent.Client,
	tursoDB *sql.DB,
	mcpClient *mcp.Client,
	history []ChatTurn,
	userMessage string,
) domainerrors.Result[*ChatAnswer] {
	if m == nil {
		return domainerrors.Err[*ChatAnswer](NewError("agent.AnswerQuery", domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}
	if tursoClient == nil || tursoDB == nil {
		return domainerrors.Err[*ChatAnswer](NewError("agent.AnswerQuery", domainerrors.CodeInvalidInput, "turso client and db cannot be nil", nil))
	}
	if mcpClient == nil {
		return domainerrors.Err[*ChatAnswer](NewError("agent.AnswerQuery", domainerrors.CodeInvalidInput, "mcp client cannot be nil", nil))
	}
	if strings.TrimSpace(userMessage) == "" {
		return domainerrors.Err[*ChatAnswer](NewError("agent.AnswerQuery", domainerrors.CodeInvalidInput, "message cannot be empty", nil))
	}

	chatTools, err := buildChatTools(tursoDB, mcpClient)
	if err != nil {
		return domainerrors.Err[*ChatAnswer](MapError("agent.AnswerQuery", "tools", err))
	}

	ag, err := llmagent.New(llmagent.Config{
		Name:                "chat_assistant",
		Description:         "Read-only operations chat assistant for the Fincher supply chain.",
		Model:               m,
		InstructionProvider: func(agent.ReadonlyContext) (string, error) { return prompts.Chat, nil },
		Tools:               chatTools,
	})
	if err != nil {
		return domainerrors.Err[*ChatAnswer](MapError("agent.AnswerQuery", "llmagent", err))
	}

	r, err := runner.NewInMemory("chat-assistant", ag)
	if err != nil {
		return domainerrors.Err[*ChatAnswer](MapError("agent.AnswerQuery", "runner", err))
	}

	prompt := renderChatHistory(history, userMessage)

	var answer string
	var citations []ChatCitation
	var runErr error
	for attempt := 1; attempt <= maxChatAttempts; attempt++ {
		answer, citations, runErr = runChatOnce(ctx, r, prompt)
		if runErr == nil {
			break
		}
		if !isTransientToolError(runErr) {
			return domainerrors.Err[*ChatAnswer](MapError("agent.AnswerQuery", "run", runErr))
		}
		logger.Warn("chat-agent: transient tool error, retrying", "attempt", attempt, "error", runErr)
	}
	if runErr != nil {
		return domainerrors.Err[*ChatAnswer](MapError("agent.AnswerQuery", "run", runErr))
	}
	if answer == "" {
		return domainerrors.Err[*ChatAnswer](NewError("agent.AnswerQuery", domainerrors.CodeInternal, "assistant produced no answer", nil))
	}

	sources := extractTitleSources(ctx, tursoClient, answer)

	logger.Info("chat-agent: answered", "citations", len(citations), "sources", len(sources), "answer_len", len(answer))

	return domainerrors.Ok(&ChatAnswer{Answer: answer, Citations: citations, Sources: sources})
}

const maxChatAttempts = 3

func isTransientToolError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unmarshal") ||
		strings.Contains(msg, "infer input schema") ||
		strings.Contains(msg, "panic in tool")
}

func runChatOnce(ctx context.Context, r *runner.Runner, prompt string) (string, []ChatCitation, error) {
	msg := genai.NewContentFromText(prompt, genai.RoleUser)

	var text strings.Builder
	citations := make([]ChatCitation, 0)
	seen := make(map[string]struct{})

	for ev, runErr := range r.Run(ctx, "operator", "chat", msg, agent.RunConfig{}) {
		if runErr != nil {
			return "", nil, runErr
		}
		if ev == nil || ev.Content == nil {
			continue
		}
		for _, p := range ev.Content.Parts {
			if p == nil {
				continue
			}
			if p.FunctionCall != nil && isRealTool(p.FunctionCall.Name) {
				label := toolCallLabel(p.FunctionCall)
				key := p.FunctionCall.Name + "|" + label
				if _, dup := seen[key]; !dup {
					seen[key] = struct{}{}
					citations = append(citations, ChatCitation{
						Label:  label,
						Source: citationSource(p.FunctionCall.Name),
						Tool:   p.FunctionCall.Name,
					})
				}
			}
			if p.Text != "" && !p.Thought {
				text.WriteString(p.Text)
			}
		}
	}

	return strings.TrimSpace(text.String()), citations, nil
}

func isRealTool(name string) bool {
	return name == "query_turso" || name == "query_clickhouse"
}

func metaSourceString(metadata map[string]any, key string) string {
	if metadata == nil {
		return ""
	}
	if v, ok := metadata[key].(string); ok {
		return v
	}
	return ""
}

func extractTitleSources(ctx context.Context, tursoClient *ent.Client, answer string) []ChatSource {
	if tursoClient == nil {
		return nil
	}
	titles, err := tursoClient.Title.Query().All(ctx)
	if err != nil {
		return nil
	}
	lowerAnswer := strings.ToLower(answer)
	withPoster := make([]ChatSource, 0)
	withoutPoster := make([]ChatSource, 0)
	seen := make(map[string]struct{})
	for _, t := range titles {
		name := strings.TrimSpace(t.Name)
		slug := strings.TrimSpace(t.Slug)
		matched := (len(name) >= 3 && strings.Contains(lowerAnswer, strings.ToLower(name))) ||
			(len(slug) >= 3 && strings.Contains(lowerAnswer, strings.ToLower(slug)))
		if !matched {
			continue
		}
		if _, dup := seen[t.ID]; dup {
			continue
		}
		seen[t.ID] = struct{}{}
		poster := metaSourceString(t.Metadata, "poster_url")
		src := ChatSource{
			Kind:        "title",
			ID:          t.ID,
			DisplayName: t.Name,
			Identifier:  t.Slug,
			Subtitle:    metaSourceString(t.Metadata, "genre"),
			ImageURL:    poster,
			Status:      string(t.OverallStatus),
		}
		if poster != "" {
			withPoster = append(withPoster, src)
		} else {
			withoutPoster = append(withoutPoster, src)
		}
	}

	// Prefer titles with poster art, cap the total to keep the payload and UI bounded.
	sources := append(withPoster, withoutPoster...)
	if len(sources) > maxChatSources {
		sources = sources[:maxChatSources]
	}
	return sources
}

const maxChatSources = 6
