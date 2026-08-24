package errors

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/logger"
)

// ErrorResponse defines the standard HTTP error payload.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Op      string `json:"op,omitempty"`
}

// Respond maps a domain error to an HTTP JSON response and logs the event.
func Respond(c echo.Context, err error) error {
	var domErr *domainerrors.DomainError
	if !errors.As(err, &domErr) || domErr == nil {
		logger.Error("internal unhandled error", "path", c.Path(), "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    string(domainerrors.CodeInternal),
			Message: "internal server error",
		})
	}

	var status int
	switch domErr.Code {
	case domainerrors.CodeNotFound:
		status = http.StatusNotFound
	case domainerrors.CodeAlreadyExists:
		status = http.StatusConflict
	case domainerrors.CodeConflict:
		status = http.StatusConflict
	case domainerrors.CodeInvalidInput:
		status = http.StatusBadRequest
	case domainerrors.CodeBudgetExceeded:
		status = http.StatusTooManyRequests
	case domainerrors.CodeUnauthenticated:
		status = http.StatusUnauthorized
	case domainerrors.CodeUnauthorized:
		status = http.StatusForbidden
	default:
		status = http.StatusInternalServerError
	}

	if status >= 500 {
		logger.Error("request failed with server error", "op", domErr.Op, "code", domErr.Code, "status", status, "error", domErr.Error())
	} else {
		logger.Warn("request failed with client error", "op", domErr.Op, "code", domErr.Code, "status", status, "message", domErr.Message)
	}

	return c.JSON(status, ErrorResponse{
		Code:    string(domErr.Code),
		Message: domErr.Message,
		Op:      domErr.Op,
	})
}
