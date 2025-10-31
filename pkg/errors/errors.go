package errors

import (
	"fmt"
)

// AppError represents application-level errors
type AppError struct {
	Code    string
	Message string
	Err     error
	Details map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New creates a new AppError
func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// Common error codes
const (
	CodeInvalidInput   = "INVALID_INPUT"
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeForbidden      = "FORBIDDEN"
	CodeNotFound       = "NOT_FOUND"
	CodeConflict       = "CONFLICT"
	CodeInternal       = "INTERNAL_ERROR"
	CodeServiceUnavail = "SERVICE_UNAVAILABLE"
	CodeTimeout        = "TIMEOUT"
	CodeBadGateway     = "BAD_GATEWAY"
	CodeTypeAssertion  = "TYPE_ASSERTION_ERROR"
	CodeParsing        = "PARSING_ERROR"
	CodeCaching        = "CACHE_ERROR"
	CodeLLMService     = "LLM_SERVICE_ERROR"
	CodeEncyclopedia   = "ENCYCLOPEDIA_SERVICE_ERROR"
)

// SafeTypeAssertion safely casts interface{} to specific types with error handling
func SafeTypeAssertion[T any](value interface{}, errMsg string) (T, error) {
	var zero T
	result, ok := value.(T)
	if !ok {
		return zero, Wrap(CodeTypeAssertion, fmt.Sprintf("type assertion failed: %s", errMsg), nil)
	}
	return result, nil
}

// SafeStringAssert safely converts interface{} to string
func SafeStringAssert(value interface{}, fieldName string) (string, error) {
	if value == nil {
		return "", New(CodeTypeAssertion, "value is nil").WithDetails(map[string]interface{}{
			"field": fieldName,
			"error": "value is nil",
		})
	}
	result, ok := value.(string)
	if !ok {
		return "", New(CodeTypeAssertion, fmt.Sprintf("invalid type for %s", fieldName)).WithDetails(map[string]interface{}{
			"field": fieldName,
			"type":  fmt.Sprintf("%T", value),
		})
	}
	return result, nil
}

// SafeFloat64Assert safely converts interface{} to float64
func SafeFloat64Assert(value interface{}, fieldName string) (float64, error) {
	if value == nil {
		return 0, New(CodeTypeAssertion, "value is nil").WithDetails(map[string]interface{}{
			"field": fieldName,
			"error": "value is nil",
		})
	}
	result, ok := value.(float64)
	if !ok {
		return 0, New(CodeTypeAssertion, fmt.Sprintf("invalid type for %s", fieldName)).WithDetails(map[string]interface{}{
			"field": fieldName,
			"type":  fmt.Sprintf("%T", value),
		})
	}
	return result, nil
}

// SafeSliceAssert safely converts interface{} to []interface{}
func SafeSliceAssert(value interface{}, fieldName string) ([]interface{}, error) {
	if value == nil {
		return []interface{}{}, New(CodeTypeAssertion, "value is nil").WithDetails(map[string]interface{}{
			"field": fieldName,
			"error": "value is nil",
		})
	}
	result, ok := value.([]interface{})
	if !ok {
		return []interface{}{}, New(CodeTypeAssertion, fmt.Sprintf("invalid type for %s", fieldName)).WithDetails(map[string]interface{}{
			"field": fieldName,
			"type":  fmt.Sprintf("%T", value),
		})
	}
	return result, nil
}

// SafeMapAssert safely converts interface{} to map[string]interface{}
func SafeMapAssert(value interface{}, fieldName string) (map[string]interface{}, error) {
	if value == nil {
		return map[string]interface{}{}, New(CodeTypeAssertion, "value is nil").WithDetails(map[string]interface{}{
			"field": fieldName,
			"error": "value is nil",
		})
	}
	result, ok := value.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, New(CodeTypeAssertion, fmt.Sprintf("invalid type for %s", fieldName)).WithDetails(map[string]interface{}{
			"field": fieldName,
			"type":  fmt.Sprintf("%T", value),
		})
	}
	return result, nil
}
