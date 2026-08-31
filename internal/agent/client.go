package agent

import (
	"context"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/gemini"
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
