package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Predefined errors
var (
	ErrUserNotFound       = errors.New("auth not found")
	ErrUserAlreadyExists  = errors.New("auth already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInternal           = errors.New("internal error")

	ErrUserAndTaskAlreadyExists = errors.New("user and task already exists")
	ErrUserAlreadyHasReferrer   = errors.New("user already has referrer")
	ErrTaskAlreadyCompleted     = errors.New("task already completed")
	ErrTaskNotFound             = errors.New("task not found")
)

type AppError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) ToHTTPResponse() (int, string) {
	return e.StatusCode, e.Message
}

func New(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func Wrap(err error, statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	switch {
	case errors.Is(err, ErrUserNotFound):
		return New(http.StatusNotFound, "User not found")
	case errors.Is(err, ErrUserAlreadyExists):
		return New(http.StatusConflict, "User already exists")
	case errors.Is(err, ErrInvalidCredentials):
		return New(http.StatusUnauthorized, "Invalid credentials")
	case errors.Is(err, ErrSessionNotFound):
		return New(http.StatusUnauthorized, "Session not found")
	case errors.Is(err, ErrInvalidInput):
		return New(http.StatusBadRequest, "Invalid input")
	case errors.Is(err, ErrUserAndTaskAlreadyExists):
		return New(http.StatusConflict, "User and task already exists")
	case errors.Is(err, ErrUserAlreadyHasReferrer):
		return New(http.StatusConflict, "User already has referrer")
	case errors.Is(err, ErrTaskAlreadyCompleted):
		return New(http.StatusConflict, "Task already completed")
	case errors.Is(err, ErrTaskNotFound):
		return New(http.StatusNotFound, "Task not found")
	default:
		return New(http.StatusInternalServerError, "Internal server error")
	}
}
