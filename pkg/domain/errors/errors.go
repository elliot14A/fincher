package errors

import (
	"fmt"
)

// ErrorCode defines canonical domain-level error classifications.
type ErrorCode string

const (
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeAlreadyExists   ErrorCode = "ALREADY_EXISTS"
	CodeInvalidInput    ErrorCode = "INVALID_INPUT"
	CodeConflict        ErrorCode = "CONFLICT"
	CodeInternal        ErrorCode = "INTERNAL"
	CodeBudgetExceeded  ErrorCode = "BUDGET_EXCEEDED"
	CodeGateTerminated  ErrorCode = "GATE_TERMINATED"
	CodeUnauthenticated ErrorCode = "UNAUTHENTICATED"
	CodeUnauthorized    ErrorCode = "UNAUTHORIZED"
)

// DomainError represents a structured, contextual error.
type DomainError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Op      string    `json:"op,omitempty"`
	Err     error     `json:"-"`
}

// New creates a DomainError.
func New(code ErrorCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewWithOp creates a DomainError with an operation name.
func NewWithOp(op string, code ErrorCode, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Op:      op,
		Err:     err,
	}
}

func (e *DomainError) Error() string {
	if e.Op != "" {
		if e.Err != nil {
			return fmt.Sprintf("[%s] %s: %s: %v", e.Code, e.Op, e.Message, e.Err)
		}
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Op, e.Message)
	}
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Option represents an optional value.
type Option[T any] struct {
	val    T
	isSome bool
}

// Some creates an Option containing a value.
func Some[T any](v T) Option[T] {
	return Option[T]{val: v, isSome: true}
}

// None creates an empty Option.
func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// IsNone returns true if the Option is empty.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// Unwrap returns the contained value, or panics if empty.
func (o Option[T]) Unwrap() T {
	if !o.isSome {
		panic("called Option.Unwrap() on a None value")
	}
	return o.val
}

// UnwrapOr returns the contained value or a fallback.
func (o Option[T]) UnwrapOr(fallback T) T {
	if !o.isSome {
		return fallback
	}
	return o.val
}

// Value returns the value and an exists flag.
func (o Option[T]) Value() (T, bool) {
	return o.val, o.isSome
}

// Result represents either success (Ok) or failure (Err).
type Result[T any] struct {
	val T
	err error
}

// Ok creates a successful Result.
func Ok[T any](val T) Result[T] {
	return Result[T]{val: val, err: nil}
}

// Err creates a failed Result.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// IsOk returns true if Result is Ok.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if Result is Err.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the value or panics if error.
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(fmt.Sprintf("called Result.Unwrap() on an Err: %v", r.err))
	}
	return r.val
}

// UnwrapOr returns the value or a fallback.
func (r Result[T]) UnwrapOr(fallback T) T {
	if r.err != nil {
		return fallback
	}
	return r.val
}

// Error returns the error if present.
func (r Result[T]) Error() error {
	return r.err
}

// Value returns the value and error tuple.
func (r Result[T]) Value() (T, error) {
	return r.val, r.err
}
