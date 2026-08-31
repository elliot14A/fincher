package agent_test

import (
	"context"
	"iter"

	"google.golang.org/adk/v2/model"
	genai "google.golang.org/genai"
)

type mockLLM struct {
	name        string
	response    string
	err         error
	lastRequest *model.LLMRequest
}

func (m *mockLLM) Name() string {
	if m.name == "" {
		return "mock-gemini"
	}
	return m.name
}

func (m *mockLLM) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	m.lastRequest = req
	return func(yield func(*model.LLMResponse, error) bool) {
		if m.err != nil {
			yield(nil, m.err)
			return
		}
		resp := &model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{Text: m.response},
				},
			},
		}
		yield(resp, nil)
	}
}
