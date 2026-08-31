package agent_test

import (
	"context"
	"testing"

	"github.com/elliot14A/fincher/internal/agent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

func TestNewModel_Validation(t *testing.T) {
	ctx := context.Background()

	t.Run("Fails when API key is empty", func(t *testing.T) {
		res := agent.NewModel(ctx, "", "gemini-2.5-flash")
		if res.IsOk() {
			t.Fatal("expected error for empty API key, got Ok")
		}
		domErr, ok := res.Error().(*domainerrors.DomainError)
		if !ok {
			t.Fatalf("expected *DomainError, got: %T", res.Error())
		}
		if domErr.Code != domainerrors.CodeInvalidInput {
			t.Errorf("expected INVALID_INPUT code, got: %s", domErr.Code)
		}
	})

	t.Run("Fails when model name is empty", func(t *testing.T) {
		res := agent.NewModel(ctx, "test-fake-key", "")
		if res.IsOk() {
			t.Fatal("expected error for empty model name, got Ok")
		}
		domErr, ok := res.Error().(*domainerrors.DomainError)
		if !ok {
			t.Fatalf("expected *DomainError, got: %T", res.Error())
		}
		if domErr.Code != domainerrors.CodeInvalidInput {
			t.Errorf("expected INVALID_INPUT code, got: %s", domErr.Code)
		}
	})

	t.Run("Initializes model with valid configuration", func(t *testing.T) {
		res := agent.NewModel(ctx, "test-fake-key", "gemini-2.5-flash")
		if res.IsErr() {
			t.Fatalf("expected successful client initialization, got error: %v", res.Error())
		}
		m := res.Unwrap()
		if m == nil {
			t.Fatal("expected non-nil model")
		}
		if m.Name() != "gemini-2.5-flash" {
			t.Errorf("expected model name gemini-2.5-flash, got: %s", m.Name())
		}
	})
}
