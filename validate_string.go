package validationcontext

import (
	"fmt"
	"os"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

var (
	emailRegexp       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	specialCharRegexp = regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`)
	numberRegexp      = regexp.MustCompile(`[0-9]`)
	uppercaseRegexp   = regexp.MustCompile(`[A-Z]`)
	lowercaseRegexp   = regexp.MustCompile(`[a-z]`)
	urlRegexp         = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
)

// ValidateMinLength checks if the value has at least minLen characters.
func (vc *ValidationContext) ValidateMinLength(value string, field string, min int, errMsg string) {
	if utf8.RuneCountInString(value) < min {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sは%d文字以上で入力してください。", field, min))
	}
}

// ValidateMaxLength checks if the value has at most maxLen characters.
func (vc *ValidationContext) ValidateMaxLength(value string, field string, max int, errMsg string) {
	if utf8.RuneCountInString(value) > max {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sは%d文字以内で入力してください。", field, max))
	}
}

// ValidateEmail checks if the value is a valid email format.
func (vc *ValidationContext) ValidateEmail(value string, field string, errMsg string) {
	if !emailRegexp.MatchString(value) {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なメールアドレスを指定してください。", field))
	}
}

// ValidateContainsSpecial checks if the value contains at least one special character.
func (vc *ValidationContext) ValidateContainsSpecial(value, field, errMsg string) {
	for _, char := range value {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			return
		}
	}
	vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、特殊文字を含めてください。", field))
}

func (vc *ValidationContext) ValidateContainsSpecialRegx(value, field, errMsg string) {
	if !specialCharRegexp.MatchString(value) {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、特殊文字を含めてください。", field))
	}
}

// ValidateContainsNumber checks if the value contains at least one number.
func (vc *ValidationContext) ValidateContainsNumber(value, field, errMsg string) {
	for _, char := range value {
		if unicode.IsDigit(char) {
			return
		}
	}
	vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、数字を含めてください。", field))
}

func (vc *ValidationContext) ValidateContainsNumberRegx(value, field, errMsg string) {
	if !numberRegexp.MatchString(value) {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、数字を含めてください。", field))
	}
}

// ValidateContainsUppercase checks if the value contains at least one uppercase letter.
func (vc *ValidationContext) ValidateContainsUppercase(value, field, errMsg string) {
	if !uppercaseRegexp.MatchString(value) {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、大文字の英字を含めてください。", field))
	}
}

// ValidateContainsLowercase checks if the value contains at least one lowercase letter.
func (vc *ValidationContext) ValidateContainsLowercase(value, field, errMsg string) {
	if !lowercaseRegexp.MatchString(value) {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、小文字の英字を含めてください。", field))
	}
}

// ValidateURL checks if the value is a valid URL.
func (vc *ValidationContext) ValidateURL(value, field, errMsg string) {
	if !urlRegexp.MatchString(value) {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なURLを指定してください。", field))
	}
}

// ValidateFile checks if the value is a valid file path.
func (vc *ValidationContext) ValidateFile(value, field, errMsg string) {
	if _, statErr := os.Stat(value); statErr != nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なファイルパスを指定してください。", field))
	}
}

// ValidateUUID checks if the value is a valid UUID.
func (vc *ValidationContext) ValidateUUID(value, field, errMsg string) {
	if _, parseErr := uuid.Parse(value); parseErr != nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なUUIDを指定してください。", field))
	}
}
