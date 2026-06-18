package validationcontext

import (
	"os"
	"strings"
	"testing"
)

func TestAddError_DefaultMessage(t *testing.T) {
	vc := NewValidationContext()
	vc.AddError("Field", "")
	if len(vc.Errors()) != 1 {
		t.Fatalf("expected 1 error, got %d", len(vc.Errors()))
	}
	if vc.Errors()[0].Message != defaultValidationMessage {
		t.Fatalf("unexpected message: %q", vc.Errors()[0].Message)
	}
}

func TestValidateFloat_DefaultMessages(t *testing.T) {
	vc := NewValidationContext()
	vc.ValidateMinFloatValue(1.0, "Score", 2.0, "")
	vc.ValidateMaxFloatValue(3.0, "Score", 2.0, "")
	vc.ValidateFloatRange(5.0, "Score", 1.0, 4.0, "")

	if len(vc.Errors()) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(vc.Errors()))
	}
	if !strings.Contains(vc.Errors()[0].Message, "Score") {
		t.Fatalf("message should include field: %q", vc.Errors()[0].Message)
	}
}

func TestValidateFile_NilSafety(t *testing.T) {
	vc := NewValidationContext()
	vc.ValidateFileExtension(nil, "Upload", []string{".png"}, "")
	vc.ValidateFileSize(nil, "Upload", 1, "")

	if len(vc.Errors()) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(vc.Errors()))
	}
}

func TestValidateFilePath_UsesDefaultWhenMissing(t *testing.T) {
	vc := NewValidationContext()
	vc.ValidateFilePath("/path/that/does/not/exist", "Path", "")
	if len(vc.Errors()) != 1 {
		t.Fatalf("expected 1 error, got %d", len(vc.Errors()))
	}
	if !strings.Contains(vc.Errors()[0].Message, "有効なファイルパス") {
		t.Fatalf("unexpected message: %s", vc.Errors()[0].Message)
	}
}

func TestValidateFileSize_WithRealFile(t *testing.T) {
	f, err := os.CreateTemp("", "validationcontext-*.txt")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if _, err := f.WriteString("1234567890"); err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}

	vc := NewValidationContext()
	vc.ValidateFileSize(f, "Upload", 1, "")
	if len(vc.Errors()) != 1 {
		t.Fatalf("expected 1 error, got %d", len(vc.Errors()))
	}
}

func TestNewValidationContext_WithStackTraceDisabled(t *testing.T) {
	vc := NewValidationContext(WithStackTrace(false))
	vc.AddError("Field", "message")
	if len(vc.Errors()) != 1 {
		t.Fatalf("expected 1 error, got %d", len(vc.Errors()))
	}
	if vc.Errors()[0].StackTrace != "" {
		t.Fatalf("expected no stack trace, got: %q", vc.Errors()[0].StackTrace)
	}
}

func TestValidateContainsRegex_BackwardCompatibility(t *testing.T) {
	vc := NewValidationContext()
	vc.ValidateContainsSpecialRegex("abc123", "FieldA", "")
	vc.ValidateContainsSpecialRegx("abc123", "FieldB", "")
	vc.ValidateContainsNumberRegex("abcdef", "FieldC", "")
	vc.ValidateContainsNumberRegx("abcdef", "FieldD", "")

	if len(vc.Errors()) != 4 {
		t.Fatalf("expected 4 errors, got %d", len(vc.Errors()))
	}
}
