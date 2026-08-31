package graph_test

import (
	"context"
	"iter"

	"google.golang.org/adk/v2/model"
	genai "google.golang.org/genai"
)

type mockLLM struct {
	name        string
	responses   []string
	callIndex   int
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

		respText := ""
		if len(m.responses) > 0 {
			if m.callIndex < len(m.responses) {
				respText = m.responses[m.callIndex]
				m.callIndex++
			} else {
				respText = m.responses[len(m.responses)-1]
			}
		}

		resp := &model.LLMResponse{
			Content: &genai.Content{
				Parts: []*genai.Part{
					{Text: respText},
				},
			},
		}
		yield(resp, nil)
	}
}
