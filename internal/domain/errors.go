package domain

import "fmt"

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeRateLimit    ErrorType = "RATE_LIMIT"
	ErrorTypeShopifyAPI   ErrorType = "SHOPIFY_API"
	ErrorTypeDatabase     ErrorType = "DATABASE"
	ErrorTypeInternal     ErrorType = "INTERNAL"
)

// AppError represents a structured application error
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a validation error
func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
	}
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(retryAfter int) *AppError {
	return &AppError{
		Type:    ErrorTypeRateLimit,
		Message: "rate limit exceeded",
		Context: map[string]interface{}{
			"retry_after": retryAfter,
		},
	}
}

// NewShopifyAPIError creates a Shopify API error
func NewShopifyAPIError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeShopifyAPI,
		Message: message,
		Err:     err,
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Err:     err,
	}
}

// NewInternalError creates an internal error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
	}
}
