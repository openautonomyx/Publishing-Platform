package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ============================================================================
// Error Types
// ============================================================================

type ErrorCode string

const (
	ErrorCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrorCodeConflict        ErrorCode = "CONFLICT"
	ErrorCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrorCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrorCodeServiceUnavail  ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
)

// AppError represents an application error with structured details
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    []string  `json:"details,omitempty"`
	RequestID  string    `json:"request_id,omitempty"`
	StatusCode int       `json:"status_code"`
	Timestamp  string    `json:"timestamp"`
}

func (e *AppError) Error() string {
	return e.Message
}

// ============================================================================
// Error Constructors
// ============================================================================

func NewValidationError(details []string) *AppError {
	return &AppError{
		Code:       ErrorCodeValidation,
		Message:    "validation failed",
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrorCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewInternalError hides the cause from the caller; log it before returning this.
func NewInternalError() *AppError {
	return &AppError{
		Code:       ErrorCodeInternalError,
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	}
}

func NewServiceUnavailableError(message string) *AppError {
	return &AppError{
		Code:       ErrorCodeServiceUnavail,
		Message:    message,
		StatusCode: http.StatusServiceUnavailable,
	}
}

func NewTooManyRequestsError() *AppError {
	return &AppError{
		Code:       ErrorCodeTooManyRequests,
		Message:    "rate limit exceeded",
		StatusCode: http.StatusTooManyRequests,
	}
}

// ============================================================================
// Error Response
// ============================================================================

type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     ErrorCode `json:"error"`
	Message   string    `json:"message"`
	Details   []string  `json:"details,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
	Timestamp string    `json:"timestamp"`
}

func NewErrorResponse(appErr *AppError, requestID string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     appErr.Code,
		Message:   appErr.Message,
		Details:   appErr.Details,
		RequestID: requestID,
		Timestamp: timeNow().UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ============================================================================
// Error Writing to HTTP Response
// ============================================================================

func writeError(w http.ResponseWriter, appErr *AppError, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)

	response := NewErrorResponse(appErr, requestID)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("failed to write error response: %v", err)
	}
}
