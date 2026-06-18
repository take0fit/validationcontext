package validationcontext

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	defaultValidationMessage = "不正な値です。"
	initialStackBufferSize   = 1024
	maxStackBufferSize       = 64 * 1024
)

type ValidationError struct {
	Field      string
	Message    string
	StackTrace string
}

type ValidationContext struct {
	errors                []ValidationError
	stackTraceEnabled bool
}

// ValidationAggregateError is a custom error type that aggregates multiple validation errors,
// including their messages and stack traces.
type ValidationAggregateError struct {
	Messages    []string
	StackTraces []string
}

// Error implements the error interface for ValidationAggregateError.
// It returns a formatted string of all validation error messages.
func (e *ValidationAggregateError) Error() string {
	return fmt.Sprintf("Validation errors: %s", strings.Join(e.Messages, "; "))
}

// GetMessages returns the list of validation error messages.
func (e *ValidationAggregateError) GetMessages() []string {
	return e.Messages
}

// GetStackTraces returns the list of stack traces associated with the validation errors.
func (e *ValidationAggregateError) GetStackTraces() []string {
	return e.StackTraces
}

// GetMessagesAsString returns all stack traces as a single string.
func (e *ValidationAggregateError) GetMessagesAsString() string {
	return strings.Join(e.Messages, "\n")
}

// GetStackTracesAsString returns all stack traces as a single string.
func (e *ValidationAggregateError) GetStackTracesAsString() string {
	return strings.Join(e.StackTraces, "\n")
}

// NewValidationContext creates and returns a new ValidationContext instance.
func NewValidationContext(opts ...ValidationContextOption) *ValidationContext {
	vc := &ValidationContext{
		errors: make([]ValidationError, 0),
		stackTraceEnabled: true,
	}
	for _, opt := range opts {
		opt(vc)
	}
	return vc
}

type ValidationContextOption func(*ValidationContext)

// WithStackTrace controls whether ValidationContext captures stack traces for each error.
func WithStackTrace(enabled bool) ValidationContextOption {
	return func(vc *ValidationContext) {
		vc.stackTraceEnabled = enabled
	}
}

// SetStackTraceEnabled updates stack trace capture behavior after initialization.
func (vc *ValidationContext) SetStackTraceEnabled(enabled bool) {
	vc.stackTraceEnabled = enabled
}

// AddError adds a validation error to the context, including the field, error message,
// and captures the stack trace at the time the error occurred.
func (vc *ValidationContext) AddError(field, message string) {
	if strings.TrimSpace(message) == "" {
		message = defaultValidationMessage
	}
	stackTrace := ""
	if vc.stackTraceEnabled {
		stackTrace = vc.captureStackTrace()
	}
	vc.errors = append(vc.errors, ValidationError{Field: field, Message: message, StackTrace: stackTrace})
}

// Errors returns the list of validation errors that have been added to the context.
func (vc *ValidationContext) Errors() []ValidationError {
	return vc.errors
}

// HasErrors returns true if there are any validation errors in the context, otherwise false.
func (vc *ValidationContext) HasErrors() bool {
	return len(vc.errors) > 0
}

// FormatErrors returns a formatted string representation of all validation errors.
func (vc *ValidationContext) FormatErrors() string {
	if !vc.HasErrors() {
		return "No validation errors"
	}
	var sb strings.Builder
	sb.WriteString("Validation errors:\n")
	for _, err := range vc.Errors() {
		sb.WriteString(fmt.Sprintf("Field: %s, Error: %s\n", err.Field, err.Message))
	}
	return sb.String()
}

// AggregateError creates and returns a ValidationAggregateError that contains
// all validation errors, including their messages and stack traces.
func (vc *ValidationContext) AggregateError() error {
	if !vc.HasErrors() {
		return nil
	}

	messages := make([]string, len(vc.errors))
	stackTraces := make([]string, len(vc.errors))

	for i, err := range vc.errors {
		messages[i] = fmt.Sprintf("Field: %s, Error: %s", err.Field, err.Message)
		stackTraces[i] = err.StackTrace
	}

	return &ValidationAggregateError{
		Messages:    messages,
		StackTraces: stackTraces,
	}
}

func (vc *ValidationContext) captureStackTrace() string {
	stackBuf := make([]byte, initialStackBufferSize)
	for {
		n := runtime.Stack(stackBuf, false)
		if n < len(stackBuf) || len(stackBuf) >= maxStackBufferSize {
			return string(stackBuf[:n])
		}
		stackBuf = make([]byte, len(stackBuf)*2)
	}
}

func (vc *ValidationContext) addErrorWithMessage(field, customMsg, defaultMsg string) {
	if strings.TrimSpace(customMsg) != "" {
		vc.AddError(field, customMsg)
		return
	}
	vc.AddError(field, defaultMsg)
}
