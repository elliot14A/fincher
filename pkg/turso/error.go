package turso

import (
	"fmt"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/ent"
)

// NewError creates a domain error.
func NewError(op string, code domainerrors.ErrorCode, msg string, err error) *domainerrors.DomainError {
	return &domainerrors.DomainError{
		Code:    code,
		Op:      "turso." + op,
		Message: msg,
		Err:     err,
	}
}

// MapEntError maps database errors to domain errors.
func MapEntError(op string, entityName string, id string, err error) error {
	if err == nil {
		return nil
	}
	if ent.IsNotFound(err) {
		return &domainerrors.DomainError{
			Code:    domainerrors.CodeNotFound,
			Op:      "turso." + op,
			Message: fmt.Sprintf("%s with id '%s' not found", entityName, id),
			Err:     err,
		}
	}
	if ent.IsConstraintError(err) {
		return &domainerrors.DomainError{
			Code:    domainerrors.CodeAlreadyExists,
			Op:      "turso." + op,
			Message: fmt.Sprintf("%s constraint violation on id '%s'", entityName, id),
			Err:     err,
		}
	}
	return &domainerrors.DomainError{
		Code:    domainerrors.CodeInternal,
		Op:      "turso." + op,
		Message: fmt.Sprintf("database error on %s", entityName),
		Err:     err,
	}
}
