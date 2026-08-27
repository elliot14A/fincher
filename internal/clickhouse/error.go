package clickhouse

import (
	"fmt"
	"strings"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
)

// NewError constructs a new domain error and logs it.
func NewError(op string, code domainerrors.ErrorCode, message string, cause error) *domainerrors.DomainError {
	domErr := domainerrors.NewWithOp(op, code, message, cause)
	if code == domainerrors.CodeInternal {
		logger.Error("clickhouse internal error", "op", op, "code", code, "message", message, "cause", cause)
	} else {
		logger.Warn("clickhouse client error", "op", op, "code", code, "message", message, "cause", cause)
	}
	return domErr
}

// MapError converts ClickHouse errors into structured domain errors and logs them.
func MapError(op, operation, identifier string, err error) *domainerrors.DomainError {
	if err == nil {
		return nil
	}

	errStr := strings.ToLower(err.Error())
	var domErr *domainerrors.DomainError

	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		domErr = domainerrors.NewWithOp(
			op,
			domainerrors.CodeInternal,
			fmt.Sprintf("clickhouse operation timed out on %s '%s'", operation, identifier),
			err,
		)
		logger.Warn("clickhouse timeout", "op", op, "operation", operation, "id", identifier)
		return domErr
	}

	if strings.Contains(errStr, "syntax error") || strings.Contains(errStr, "type mismatch") || strings.Contains(errStr, "bad_arguments") {
		domErr = domainerrors.NewWithOp(
			op,
			domainerrors.CodeInvalidInput,
			fmt.Sprintf("invalid query or parameters for %s '%s'", operation, identifier),
			err,
		)
		logger.Warn("clickhouse query error", "op", op, "operation", operation, "id", identifier, "cause", err)
		return domErr
	}

	domErr = domainerrors.NewWithOp(
		op,
		domainerrors.CodeInternal,
		fmt.Sprintf("clickhouse failure for %s '%s'", operation, identifier),
		err,
	)
	logger.Error("clickhouse failure", "op", op, "operation", operation, "id", identifier, "cause", err)
	return domErr
}
