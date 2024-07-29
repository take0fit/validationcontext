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
		name   string
		errors []ValidationError
		want   bool
	}{
		{"No errors", []ValidationError{}, false},
		{"With errors", []ValidationError{{Field: "Field1", Message: "Error1"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &ValidationContext{errors: tt.errors}
			if got := vc.HasErrors(); got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatErrors(t *testing.T) {
	tests := []struct {
		name   string
		errors []ValidationError
		want   string
	}{
		{"No errors", []ValidationError{}, "No validation errors"},
		{"With errors", []ValidationError{{Field: "Field1", Message: "Error1"}, {Field: "Field2", Message: "Error2"}}, "Validation errors:\nField: Field1, Error: Error1\nField: Field2, Error: Error2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &ValidationContext{errors: tt.errors}
			if got := vc.FormatErrors(); got != tt.want {
				t.Errorf("FormatErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregateError(t *testing.T) {
	tests := []struct {
		name   string
		errors []ValidationError
		want   string
	}{
		{"No errors", []ValidationError{}, ""},
		{"With errors", []ValidationError{{Field: "Field1", Message: "Error1"}, {Field: "Field2", Message: "Error2"}}, "Validation errors:\nField: Field1, Error: Error1\nField: Field2, Error: Error2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := &ValidationContext{errors: tt.errors}
			err := vc.AggregateError()
			if (err == nil && tt.want != "") || (err != nil && err.Error() != tt.want) {
				t.Errorf("AggregateError() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestValidateMinLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		field     string
		minLen    int
		errMsg    string
		wantError bool
	}{
		{"Min length not met", "abc", "MinLengthField", 5, "Min length is 5", true},
		{"Min length met", "abcdef", "MinLengthField", 5, "Min length is 5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMinLength(tt.value, tt.field, tt.minLen, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMinLength() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateMaxLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		field     string
		maxLen    int
		errMsg    string
		wantError bool
	}{
		{"Max length exceeded", "abcdefghij", "MaxLengthField", 5, "Max length is 5", true},
		{"Max length not exceeded", "abc", "MaxLengthField", 5, "Max length is 5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMaxLength(tt.value, tt.field, tt.maxLen, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMaxLength() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateEmailFormat(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		field     string
		errMsg    string
		wantError bool
	}{
		{"Invalid email format", "invalid-email", "EmailField", "Invalid email format", true},
		{"Valid email format", "test@example.com", "EmailField", "Invalid email format", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateEmail(tt.email, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateEmailFormat() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateContainsSpecial(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		field     string
		errMsg    string
		wantError bool
	}{
		{"No special character", "abc123", "SpecialCharField", "Must contain special character", true},
		{"With special character", "abc!123", "SpecialCharField", "Must contain special character", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsSpecial(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsSpecial() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateContainsNumber(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		field     string
		errMsg    string
		wantError bool
	}{
		{"No number", "abcdef", "NumberField", "Must contain number", true},
		{"With number", "abc123", "NumberField", "Must contain number", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsNumber(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsNumber() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateContainsUppercase(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		field     string
		errMsg    string
		wantError bool
	}{
		{"No uppercase letter", "abcdef", "UppercaseField", "Must contain uppercase letter", true},
		{"With uppercase letter", "Abcdef", "UppercaseField", "Must contain uppercase letter", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsUppercase(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsUppercase() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateContainsLowercase(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		field     string
		errMsg    string
		wantError bool
	}{
		{"No lowercase letter", "ABCDEF", "LowercaseField", "Must contain lowercase letter", true},
		{"With lowercase letter", "Abcdef", "LowercaseField", "Must contain lowercase letter", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateContainsLowercase(tt.value, tt.field, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateContainsLowercase() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateMinValue(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		field     string
		minValue  int
		errMsg    string
		wantError bool
	}{
		{"Value below min", 5, "MinValueField", 10, "Min value is 10", true},
		{"Value above min", 15, "MinValueField", 10, "Min value is 10", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMinValue(tt.value, tt.field, tt.minValue, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMinValue() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateMaxValue(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		field     string
		maxValue  int
		errMsg    string
		wantError bool
	}{
		{"Value above max", 15, "MaxValueField", 10, "Max value is 10", true},
		{"Value below max", 5, "MaxValueField", 10, "Max value is 10", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMaxValue(tt.value, tt.field, tt.maxValue, tt.errMsg)
			if vc.HasErrors() != tt.wantError {
				t.Errorf("ValidateMaxValue() = %v, want %v", vc.HasErrors(), tt.wantError)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidDate", "2023-07-25", false},
		{"InvalidDate", "2023-07-32", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateDate(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestValidateYearMonth(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidYearMonth", "2023-07", false},
		{"InvalidYearMonth", "2023-13", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateYearMonth(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestValidateYear(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidYear", "2023", false},
		{"InvalidYear", "20A3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateYear(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestValidateMonth(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidMonth", "07", false},
		{"InvalidMonth", "13", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateMonth(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestValidateDateTime(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidDateTime", "2023-07-25 15:04:05", false},
		{"InvalidDateTime", "2023-07-25 25:04:05", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateDateTime(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestValidateTime(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidTime", "15:04", false},
		{"InvalidTime", "25:04", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateTime(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidURL", "https://www.example.com", false},
		{"InvalidURL", "ht://www.example.com", true},
		{"EmptyScheme", "www.example.com", true},
		{"EmptyHost", "https://", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateURL(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
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
		name      string
		value     string
		expectErr bool
	}{
		{"ValidFile", tmpFile.Name(), false},
		{"InvalidFile", "nonexistentfile.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateFilePath(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
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
		expectErr       bool
	}{
		{"ValidExtension", tmpFile, []string{".txt", ".md"}, false},
		{"InvalidExtension", tmpFile, []string{".exe", ".bin"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateFileExtension(tt.file, "Field1", tt.validExtensions, "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"ValidUUID", "123e4567-e89b-12d3-a456-426614174000", false},
		{"InvalidUUID", "invalid-uuid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.ValidateUUID(tt.value, "Field1", "")
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		skipNil   bool
		expectErr bool
	}{
		{"NonEmptyString", "test", false, false},
		{"EmptyString", "", false, true},
		{"NilPointer", (*string)(nil), false, true},
		{"ZeroInt", 0, false, true},
		{"NonZeroInt", 123, false, false},
		{"EmptySlice", []int{}, false, true},
		{"NonEmptySlice", []int{1, 2, 3}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := NewValidationContext()
			vc.Required(tt.value, "Field1", "", tt.skipNil)
			if (len(vc.Errors()) > 0) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, vc.Errors())
			}
		})
	}
}
