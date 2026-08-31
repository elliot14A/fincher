package agent

import (
	"fmt"
	"strings"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
)

// NewError constructs a structured domain error for the agent service and logs it.
func NewError(op string, code domainerrors.ErrorCode, message string, cause error) *domainerrors.DomainError {
	domErr := domainerrors.NewWithOp(op, code, message, cause)
	if code == domainerrors.CodeInternal {
		logger.Error("agent internal error", "op", op, "code", code, "message", message, "cause", cause)
	} else {
		logger.Warn("agent client error", "op", op, "code", code, "message", message, "cause", cause)
	}
	return domErr
}

// MapError converts agent, model, and workflow errors into structured domain errors.
func MapError(op, target string, err error) *domainerrors.DomainError {
	if err == nil {
		return nil
	}

	errStr := strings.ToLower(err.Error())
	var domErr *domainerrors.DomainError

	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		domErr = domainerrors.NewWithOp(
			op,
			domainerrors.CodeInternal,
			fmt.Sprintf("agent operation timed out on %s", target),
			err,
		)
		logger.Warn("agent timeout", "op", op, "target", target, "cause", err)
		return domErr
	}

	if strings.Contains(errStr, "api key") || strings.Contains(errStr, "unauthenticated") {
		domErr = domainerrors.NewWithOp(
			op,
			domainerrors.CodeUnauthenticated,
			fmt.Sprintf("agent authentication failure on %s", target),
			err,
		)
		logger.Error("agent authentication error", "op", op, "target", target, "cause", err)
		return domErr
	}

	if strings.Contains(errStr, "quota") || strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "resource exhausted") {
		domErr = domainerrors.NewWithOp(
			op,
			domainerrors.CodeBudgetExceeded,
			fmt.Sprintf("model budget/quota exceeded on %s", target),
			err,
		)
		logger.Warn("model quota exceeded", "op", op, "target", target, "cause", err)
		return domErr
	}

	domErr = domainerrors.NewWithOp(
		op,
		domainerrors.CodeInternal,
		fmt.Sprintf("agent execution failure on %s", target),
		err,
	)
	logger.Error("agent failure", "op", op, "target", target, "cause", err)
	return domErr
}
