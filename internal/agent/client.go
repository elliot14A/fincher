package agent

import (
	"context"
	"database/sql"

	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/tool"
	genai "google.golang.org/genai"
)

// NewModel initializes an ADK model.LLM backed by Google Gemini.
func NewModel(ctx context.Context, apiKey, modelName string) domainerrors.Result[model.LLM] {
	if apiKey == "" {
		return domainerrors.Err[model.LLM](NewError("agent.NewModel", domainerrors.CodeInvalidInput, "gemini api key is required", nil))
	}
	if modelName == "" {
		return domainerrors.Err[model.LLM](NewError("agent.NewModel", domainerrors.CodeInvalidInput, "model name is required", nil))
	}

	cfg := &genai.ClientConfig{
		APIKey: apiKey,
	}

	m, err := gemini.NewModel(ctx, modelName, cfg)
	if err != nil {
		return domainerrors.Err[model.LLM](MapError("agent.NewModel", modelName, err))
	}

	return domainerrors.Ok(m)
}

// BuildAgentTools constructs the complete ADK toolset for AI evaluation agents.
func BuildAgentTools(tursoClient *ent.Client, chDB *sql.DB) ([]tool.Tool, error) {
	if tursoClient == nil {
		return nil, domainerrors.NewWithOp("agent.BuildAgentTools", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	if chDB == nil {
		return nil, domainerrors.NewWithOp("agent.BuildAgentTools", domainerrors.CodeInvalidInput, "clickhouse db cannot be nil", nil)
	}

	analyticsTool, err := tools.NewAnalyticsTool(chDB)
	if err != nil {
		return nil, err
	}
	impactTool, err := tools.NewDeliveryImpactTool(tursoClient)
	if err != nil {
		return nil, err
	}
	candidatesTool, err := tools.NewVendorCandidatesTool(tursoClient, chDB)
	if err != nil {
		return nil, err
	}
	projectionTool, err := tools.NewProjectionTool(tursoClient)
	if err != nil {
		return nil, err
	}

	return []tool.Tool{analyticsTool, impactTool, candidatesTool, projectionTool}, nil
}
