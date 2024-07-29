package validationcontext

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

type ValidationContext struct {
	errors []ValidationError
}

func NewValidationContext() *ValidationContext {
	return &ValidationContext{
		errors: make([]ValidationError, 0),
	}
}

func (vc *ValidationContext) AddError(field, message string) {
	vc.errors = append(vc.errors, ValidationError{Field: field, Message: message})
}

func (vc *ValidationContext) Errors() []ValidationError {
	return vc.errors
}

func (vc *ValidationContext) HasErrors() bool {
	return len(vc.errors) > 0
}

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

func (vc *ValidationContext) AggregateError() error {
	if !vc.HasErrors() {
		return nil
	}
	return fmt.Errorf(vc.FormatErrors())
}
