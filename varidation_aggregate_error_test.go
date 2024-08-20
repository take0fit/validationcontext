package validationcontext

import (
	"testing"
)

func TestValidationAggregateError_Error(t *testing.T) {
	aggErr := &ValidationAggregateError{
		Messages: []string{"Error1", "Error2"},
	}

	expected := "Validation errors: Error1; Error2"
	if aggErr.Error() != expected {
		t.Errorf("Expected: %s, got: %s", expected, aggErr.Error())
	}
}

func TestValidationAggregateError_GetMessages(t *testing.T) {
	aggErr := &ValidationAggregateError{
		Messages: []string{"Error1", "Error2"},
	}

	messages := aggErr.GetMessages()
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}
	if messages[0] != "Error1" || messages[1] != "Error2" {
		t.Errorf("Unexpected messages: %v", messages)
	}
}

func TestValidationAggregateError_GetStackTraces(t *testing.T) {
	aggErr := &ValidationAggregateError{
		StackTraces: []string{"Trace1", "Trace2"},
	}

	stackTraces := aggErr.GetStackTraces()
	if len(stackTraces) != 2 {
		t.Errorf("Expected 2 stack traces, got %d", len(stackTraces))
	}
	if stackTraces[0] != "Trace1" || stackTraces[1] != "Trace2" {
		t.Errorf("Unexpected stack traces: %v", stackTraces)
	}
}
