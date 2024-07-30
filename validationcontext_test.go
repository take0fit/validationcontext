package validationcontext

import (
	"os"
	"testing"
)

func TestAddError(t *testing.T) {
	vc := NewValidationContext()
	vc.AddError("Field1", "Error1")
	if len(vc.Errors()) != 1 {
		t.Errorf("Expected 1 error, got %d", len(vc.Errors()))
	}
	if vc.Errors()[0].Field != "Field1" || vc.Errors()[0].Message != "Error1" {
		t.Errorf("Unexpected error: %v", vc.Errors()[0])
	}
}

func TestErrors(t *testing.T) {
	vc := NewValidationContext()
	vc.AddError("Field1", "Error1")
	vc.AddError("Field2", "Error2")
	if len(vc.Errors()) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(vc.Errors()))
	}
}

func TestHasErrors(t *testing.T) {
	tests := []struct {
		name           string
		errors         []ValidationError
		want           bool
		expectErrCount int
	}{
		{"No errors", []ValidationError{}, false, 0},
		{"With errors", []ValidationError{{Field: "Field1", Message: "Error1"}}, true, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &ValidationContext{errors: tt.errors}
			if got := vc.HasErrors(); got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestFormatErrors(t *testing.T) {
	tests := []struct {
		name           string
		errors         []ValidationError
		want           string
		expectErrCount int
	}{
		{"No errors", []ValidationError{}, "No validation errors", 0},
		{"With errors", []ValidationError{{Field: "Field1", Message: "Error1"}, {Field: "Field2", Message: "Error2"}}, "Validation errors:\nField: Field1, Error: Error1\nField: Field2, Error: Error2\n", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &ValidationContext{errors: tt.errors}
			if got := vc.FormatErrors(); got != tt.want {
				t.Errorf("FormatErrors() = %v, want %v", got, tt.want)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestAggregateError(t *testing.T) {
	tests := []struct {
		name           string
		errors         []ValidationError
		want           string
		expectErrCount int
	}{
		{"No errors", []ValidationError{}, "", 0},
		{"With errors", []ValidationError{{Field: "Field1", Message: "Error1"}, {Field: "Field2", Message: "Error2"}}, "Validation errors:\nField: Field1, Error: Error1\nField: Field2, Error: Error2\n", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &ValidationContext{errors: tt.errors}
			err := vc.AggregateError()
			if (err == nil && tt.want != "") || (err != nil && err.Error() != tt.want) {
				t.Errorf("AggregateError() = %v, want %v", err, tt.want)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateMinLength(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		field          string
		minLen         int
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"Min length not met", "abc", "MinLengthField", 5, "Min length is 5", true, 1},
		{"Min length met", "abcdef", "MinLengthField", 5, "Min length is 5", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMinLength(tt.value, tt.field, tt.minLen, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMinLength() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateMaxLength(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		field          string
		maxLen         int
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"Max length exceeded", "abcdefghij", "MaxLengthField", 5, "Max length is 5", true, 1},
		{"Max length not exceeded", "abc", "MaxLengthField", 5, "Max length is 5", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMaxLength(tt.value, tt.field, tt.maxLen, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMaxLength() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateEmailFormat(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		field          string
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"Invalid email format", "invalid-email", "EmailField", "Invalid email format", true, 1},
		{"Valid email format", "test@example.com", "EmailField", "Invalid email format", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateEmail(tt.email, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateEmailFormat() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateContainsSpecial(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		field          string
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"No special character", "abc123", "SpecialCharField", "Must contain special character", true, 1},
		{"With special character", "abc!123", "SpecialCharField", "Must contain special character", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsSpecial(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsSpecial() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateContainsNumber(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		field          string
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"No number", "abcdef", "NumberField", "Must contain number", true, 1},
		{"With number", "abc123", "NumberField", "Must contain number", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsNumber(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsNumber() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateContainsUppercase(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		field          string
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"No uppercase letter", "abcdef", "UppercaseField", "Must contain uppercase letter", true, 1},
		{"With uppercase letter", "Abcdef", "UppercaseField", "Must contain uppercase letter", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsUppercase(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsUppercase() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateContainsLowercase(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		field          string
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"No lowercase letter", "ABCDEF", "LowercaseField", "Must contain lowercase letter", true, 1},
		{"With lowercase letter", "Abcdef", "LowercaseField", "Must contain lowercase letter", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsLowercase(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsLowercase() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateMinValue(t *testing.T) {
	tests := []struct {
		name           string
		value          int
		field          string
		minValue       int
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"Value below min", 5, "MinValueField", 10, "Min value is 10", true, 1},
		{"Value above min", 15, "MinValueField", 10, "Min value is 10", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMinValue(tt.value, tt.field, tt.minValue, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMinValue() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateMaxValue(t *testing.T) {
	tests := []struct {
		name           string
		value          int
		field          string
		maxValue       int
		errMsg         string
		wantError      bool
		expectErrCount int
	}{
		{"Value above max", 15, "MaxValueField", 10, "Max value is 10", true, 1},
		{"Value below max", 5, "MaxValueField", 10, "Max value is 10", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMaxValue(tt.value, tt.field, tt.maxValue, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMaxValue() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidDate", "2023-07-25", 0},
		{"InvalidDate", "2023-07-32", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateDate(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateYearMonth(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidYearMonth", "2023-07", 0},
		{"InvalidYearMonth", "2023-13", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateYearMonth(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateYear(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidYear", "2023", 0},
		{"InvalidYear", "20A3", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateYear(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateMonth(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidMonth", "07", 0},
		{"InvalidMonth", "13", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMonth(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateDateTime(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidDateTime", "2023-07-25 15:04:05", 0},
		{"InvalidDateTime", "2023-07-25 25:04:05", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateDateTime(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateTime(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidTime", "15:04", 0},
		{"InvalidTime", "25:04", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateTime(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidURL", "https://www.example.com", 0},
		{"InvalidURL", "ht://www.example.com", 1},
		{"EmptyScheme", "www.example.com", 1},
		{"EmptyHost", "https://", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateURL(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateFile(t *testing.T) {
	// テスト用の一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidFile", tmpFile.Name(), 0},
		{"InvalidFile", "nonexistentfile.txt", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateFilePath(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateFileExtension(t *testing.T) {
	// テスト用の一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tests := []struct {
		name            string
		file            *os.File
		validExtensions []string
		expectErrCount  int
	}{
		{"ValidExtension", tmpFile, []string{".txt", ".md"}, 0},
		{"InvalidExtension", tmpFile, []string{".exe", ".bin"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateFileExtension(tt.file, "Field1", tt.validExtensions, "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		expectErrCount int
	}{
		{"ValidUUID", "123e4567-e89b-12d3-a456-426614174000", 0},
		{"InvalidUUID", "invalid-uuid", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateUUID(tt.value, "Field1", "")
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		skipNil        bool
		expectErrCount int
	}{
		{"NonEmptyString", "test", false, 0},
		{"EmptyString", "", false, 1},
		{"NilPointer", (*string)(nil), false, 1},
		{"ZeroInt", 0, false, 1},
		{"NonZeroInt", 123, false, 0},
		{"EmptySlice", []int{}, false, 1},
		{"NonEmptySlice", []int{1, 2, 3}, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.Required(tt.value, "Field1", "", tt.skipNil)
			if len(vc.Errors()) != tt.expectErrCount {
				t.Errorf("Expected error count: %v, got: %v", tt.expectErrCount, len(vc.Errors()))
			}
		})
	}
}
