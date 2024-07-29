package validationcontext

import (
	"regexp"
	"unicode"
)

// ValidateRequired checks if the value is not empty.
func (vc *ValidationContext) ValidateRequired(value, field, errMsg string) {
	if value == "" {
		vc.AddError(field, errMsg)
	}
}

// ValidateMinLength checks if the value has at least minLen characters.
func (vc *ValidationContext) ValidateMinLength(value, field string, minLen int, errMsg string) {
	if len(value) < minLen {
		vc.AddError(field, errMsg)
	}
}

// ValidateMaxLength checks if the value has at most maxLen characters.
func (vc *ValidationContext) ValidateMaxLength(value, field string, maxLen int, errMsg string) {
	if len(value) > maxLen {
		vc.AddError(field, errMsg)
	}
}

// ValidateEmailFormat checks if the value is a valid email format.
func (vc *ValidationContext) ValidateEmailFormat(email, field, errMsg string) {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		vc.AddError(field, errMsg)
	}
}

// ValidateContainsSpecial checks if the value contains at least one special character.
func (vc *ValidationContext) ValidateContainsSpecial(value, field, errMsg string) {
	hasSpecial := false
	for _, char := range value {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		vc.AddError(field, errMsg)
	}
}

// ValidateContainsNumber checks if the value contains at least one number.
func (vc *ValidationContext) ValidateContainsNumber(value, field, errMsg string) {
	hasNumber := false
	for _, char := range value {
		if unicode.IsDigit(char) {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		vc.AddError(field, errMsg)
	}
}

// ValidateContainsUppercase checks if the value contains at least one uppercase letter.
func (vc *ValidationContext) ValidateContainsUppercase(value, field, errMsg string) {
	hasUpper := false
	for _, char := range value {
		if unicode.IsUpper(char) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		vc.AddError(field, errMsg)
	}
}

// ValidateContainsLowercase checks if the value contains at least one lowercase letter.
func (vc *ValidationContext) ValidateContainsLowercase(value, field, errMsg string) {
	hasLower := false
	for _, char := range value {
		if unicode.IsLower(char) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		vc.AddError(field, errMsg)
	}
}
