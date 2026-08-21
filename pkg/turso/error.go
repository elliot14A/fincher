package turso

import (
	"fmt"
	"log/slog"
	"strings"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/ent"
)

// NewError constructs a new domain error and logs it.
func NewError(op string, code domainerrors.ErrorCode, message string, cause error) *domainerrors.DomainError {
	domErr := domainerrors.NewWithOp(op, code, message, cause)
	slog.Error("turso store error", "op", op, "code", code, "message", message, "cause", cause)
	return domErr
}

// MapEntError converts Ent/database errors into domain errors and logs them.
func MapEntError(op, entity, id string, err error) *domainerrors.DomainError {
	if err == nil {
		return nil
	}

	var domErr *domainerrors.DomainError

	if ent.IsNotFound(err) {
		domErr = domainerrors.NewWithOp(
			op,
			domainerrors.CodeNotFound,
			fmt.Sprintf("%s with id '%s' not found", entity, id),
			err,
		)
		slog.Warn("record not found", "op", op, "entity", entity, "id", id)
		return domErr
	}

	if ent.IsConstraintError(err) {
		errStr := strings.ToLower(err.Error())
		isDelete := strings.Contains(strings.ToLower(op), "delete")

		if strings.Contains(errStr, "foreign key") {
			if isDelete {
				domErr = domainerrors.NewWithOp(
					op,
					domainerrors.CodeConflict,
					fmt.Sprintf("%s '%s' cannot be deleted because dependent records reference it", entity, id),
					err,
				)
			} else {
				domErr = domainerrors.NewWithOp(
					op,
					domainerrors.CodeInvalidInput,
					fmt.Sprintf("%s foreign key reference invalid on id '%s'", entity, id),
					err,
				)
			}
		} else {
			domErr = domainerrors.NewWithOp(
				op,
				domainerrors.CodeAlreadyExists,
				fmt.Sprintf("%s constraint violation on id '%s'", entity, id),
				err,
			)
		}
		slog.Warn("constraint error", "op", op, "entity", entity, "id", id, "code", domErr.Code, "cause", err)
		return domErr
	}

	domErr = domainerrors.NewWithOp(
		op,
		domainerrors.CodeInternal,
		fmt.Sprintf("database failure for %s '%s'", entity, id),
		err,
	)
	slog.Error("database failure", "op", op, "entity", entity, "id", id, "cause", err)
	return domErr
}
