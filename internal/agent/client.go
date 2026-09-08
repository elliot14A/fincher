package agent

import (
	"context"
	"fmt"

	"github.com/elliot14A/fincher/internal/agent/tools"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/mcp"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/tool"
	genai "google.golang.org/genai"
)

func NewModel(ctx context.Context, apiKey, modelName string, opts map[string]string) domainerrors.Result[model.LLM] {
	if apiKey == "" {
		return domainerrors.Err[model.LLM](NewError("agent.NewModel", domainerrors.CodeInvalidInput, "gemini api key is required", nil))
	}
	if modelName == "" {
		return domainerrors.Err[model.LLM](NewError("agent.NewModel", domainerrors.CodeInvalidInput, "model name is required", nil))
	}

	cfg := &genai.ClientConfig{
		Backend: genai.BackendVertexAI,
		APIKey:  apiKey,
	}

	if loc := opts["location"]; loc != "" {
		cfg.HTTPOptions.BaseURL = fmt.Sprintf("https://%s-aiplatform.googleapis.com/", loc)
	}

	m, err := gemini.NewModel(ctx, modelName, cfg)
	if err != nil {
		return domainerrors.Err[model.LLM](MapError("agent.NewModel", modelName, err))
	}

	return domainerrors.Ok(m)
}

func BuildAgentTools(tursoClient *ent.Client, mcpClient *mcp.Client) ([]tool.Tool, error) {
	if tursoClient == nil {
		return nil, domainerrors.NewWithOp("agent.BuildAgentTools", domainerrors.CodeInvalidInput, "turso client cannot be nil", nil)
	}
	if mcpClient == nil {
		return nil, domainerrors.NewWithOp("agent.BuildAgentTools", domainerrors.CodeInvalidInput, "mcp client cannot be nil", nil)
	}

	analyticsTool, err := tools.NewAnalyticsTool(mcpClient)
	if err != nil {
		return nil, err
	}
	impactTool, err := tools.NewDeliveryImpactTool(tursoClient)
	if err != nil {
		return nil, err
	}
	candidatesTool, err := tools.NewVendorCandidatesTool(tursoClient, mcpClient)
	if err != nil {
		return nil, err
	}
	projectionTool, err := tools.NewProjectionTool(tursoClient)
	if err != nil {
		return nil, err
	}

	return []tool.Tool{analyticsTool, impactTool, candidatesTool, projectionTool}, nil
}
