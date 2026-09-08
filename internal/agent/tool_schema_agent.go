package agent

import (
	"context"
	"strings"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/genai"
)

func runToolPhase(ctx context.Context, m model.LLM, appName, sessionID, instruction, userPrompt string, tools []tool.Tool) (string, int, error) {
	ag, err := llmagent.New(llmagent.Config{
		Name:                appName + "_gather",
		Description:         "Gathers operational data via read-only queries.",
		Model:               m,
		InstructionProvider: func(agent.ReadonlyContext) (string, error) { return instruction, nil },
		Tools:               tools,
	})
	if err != nil {
		return "", 0, err
	}
	r, err := runner.NewInMemory(appName+"-gather", ag)
	if err != nil {
		return "", 0, err
	}

	msg := genai.NewContentFromText(userPrompt, genai.RoleUser)
	var text strings.Builder
	toolCalls := 0
	for ev, runErr := range r.Run(ctx, "system", sessionID, msg, agent.RunConfig{}) {
		if runErr != nil {
			return "", toolCalls, runErr
		}
		if ev == nil || ev.Content == nil {
			continue
		}
		for _, p := range ev.Content.Parts {
			if p == nil {
				continue
			}
			if p.FunctionCall != nil {
				toolCalls++
			}
			if p.Text != "" && !p.Thought {
				text.WriteString(p.Text)
				text.WriteString("\n")
			}
		}
	}
	return strings.TrimSpace(text.String()), toolCalls, nil
}

func runSchemaPhase(ctx context.Context, m model.LLM, appName, sessionID, instruction, userPrompt string, schema *genai.Schema) (string, error) {
	outKey := appName + "_out"
	ag, err := llmagent.New(llmagent.Config{
		Name:                appName + "_emit",
		Description:         "Emits the final structured decision.",
		Model:               m,
		InstructionProvider: func(agent.ReadonlyContext) (string, error) { return instruction, nil },
		OutputSchema:        schema,
		OutputKey:           outKey,
	})
	if err != nil {
		return "", err
	}
	r, err := runner.NewInMemory(appName+"-emit", ag)
	if err != nil {
		return "", err
	}

	msg := genai.NewContentFromText(userPrompt, genai.RoleUser)
	var rawJSON, lastText string
	for ev, runErr := range r.Run(ctx, "system", sessionID, msg, agent.RunConfig{}) {
		if runErr != nil {
			return "", runErr
		}
		if ev == nil {
			continue
		}
		if ev.Actions.StateDelta != nil {
			if v, ok := ev.Actions.StateDelta[outKey]; ok {
				if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
					rawJSON = s
				}
			}
		}
		if ev.Content != nil {
			for _, p := range ev.Content.Parts {
				if p != nil && p.Text != "" && !p.Thought {
					lastText = p.Text
				}
			}
		}
	}
	if rawJSON == "" {
		rawJSON = lastText
	}
	clean := strings.TrimSpace(rawJSON)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)
	if clean == "" {
		return "", domainerrors.NewWithOp("agent.runSchemaPhase", domainerrors.CodeInternal, "schema phase produced no output", nil)
	}
	return clean, nil
}

func runToolThenSchema(ctx context.Context, m model.LLM, appName, sessionID, gatherInstruction, gatherPrompt string, tools []tool.Tool, emitInstruction string, schema *genai.Schema) (string, int, error) {
	findings, toolCalls, err := runToolPhase(ctx, m, appName, sessionID, gatherInstruction, gatherPrompt, tools)
	if err != nil {
		return "", toolCalls, err
	}
	logger.Info("tool-schema-agent: gather phase complete", "app", appName, "tool_calls", toolCalls, "findings_len", len(findings))

	emitPrompt := gatherPrompt + "\n\nYour research and vendor findings:\n" + findings + "\n\nNow produce the final structured output based strictly on the findings above."
	structured, err := runSchemaPhase(ctx, m, appName, sessionID, emitInstruction, emitPrompt, schema)
	if err != nil {
		return "", toolCalls, err
	}
	return structured, toolCalls, nil
}
