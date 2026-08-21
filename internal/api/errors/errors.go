package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// ErrorResponse defines standard JSON error response body.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Respond maps domain error to HTTP status code and JSON response.
func Respond(c echo.Context, err error) error {
	var domErr *domainerrors.DomainError
	if errors.As(err, &domErr) {
		status := http.StatusInternalServerError
		switch domErr.Code {
		case domainerrors.CodeNotFound:
			status = http.StatusNotFound
		case domainerrors.CodeInvalidInput:
			status = http.StatusBadRequest
		case domainerrors.CodeAlreadyExists, domainerrors.CodeConflict:
			status = http.StatusConflict
		case domainerrors.CodeUnauthenticated:
			status = http.StatusUnauthorized
		case domainerrors.CodeUnauthorized:
			status = http.StatusForbidden
		case domainerrors.CodeBudgetExceeded:
			status = http.StatusTooManyRequests
		}
		return c.JSON(status, ErrorResponse{
			Code:    string(domErr.Code),
			Message: domErr.Message,
		})
	}

	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    string(domainerrors.CodeInternal),
		Message: err.Error(),
	})
}
