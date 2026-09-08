package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/adk/v2/model"
	genai "google.golang.org/genai"
)

func generateJSON[T any](ctx context.Context, m model.LLM, op, systemPrompt, userPrompt string) domainerrors.Result[T] {
	if m == nil {
		return domainerrors.Err[T](NewError(op, domainerrors.CodeInvalidInput, "llm model cannot be nil", nil))
	}

	noThinking := int32(0)
	genConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: systemPrompt},
			},
		},
		ThinkingConfig: &genai.ThinkingConfig{ThinkingBudget: &noThinking},
	}
	if schema, err := jsonschema.For[T](nil); err == nil && schema != nil {
		genConfig.ResponseJsonSchema = schema
	}

	req := &model.LLMRequest{
		Model: m.Name(),
		Contents: []*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{Text: userPrompt},
				},
			},
		},
		Config: genConfig,
	}

	var rawResponse strings.Builder
	for resp, err := range m.GenerateContent(ctx, req, false) {
		if err != nil {
			return domainerrors.Err[T](MapError(op, m.Name(), err))
		}
		if resp != nil && resp.Content != nil {
			for _, part := range resp.Content.Parts {
				if part != nil && part.Text != "" {
					rawResponse.WriteString(part.Text)
				}
			}
		}
	}

	rawText := strings.TrimSpace(rawResponse.String())
	if rawText == "" {
		return domainerrors.Err[T](NewError(op, domainerrors.CodeInternal, "model returned empty response", nil))
	}

	cleanJSON := strings.TrimPrefix(rawText, "```json")
	cleanJSON = strings.TrimPrefix(cleanJSON, "```")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	var target T
	if err := json.Unmarshal([]byte(cleanJSON), &target); err != nil {
		return domainerrors.Err[T](NewError(op, domainerrors.CodeInternal, fmt.Sprintf("failed to parse json response: %s", cleanJSON), err))
	}

	return domainerrors.Ok(target)
}
