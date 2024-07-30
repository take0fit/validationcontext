package validationcontext

import (
	"fmt"
	"os"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ValidateMinLength checks if the value has at least minLen characters.
func (vc *ValidationContext) ValidateMinLength(value string, field string, min int, errMsg string) {
	if utf8.RuneCountInString(value) < min {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sは%d文字以上で入力してください。", field, min))
	}
}

// ValidateMaxLength checks if the value has at most maxLen characters.
func (vc *ValidationContext) ValidateMaxLength(value string, field string, max int, errMsg string) {
	if utf8.RuneCountInString(value) > max {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sは%d文字以内で入力してください。", field, max))
	}
}

// ValidateEmail checks if the value is a valid email format.
func (vc *ValidationContext) ValidateEmail(value string, field string, errMsg string) {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(value) {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効なメールアドレスを指定してください。", field))
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
	if hasSpecial {
		return
	}
	if errMsg != "" {
		vc.AddError(field, errMsg)
		return
	}
	vc.AddError(field, fmt.Sprintf("%sには、特殊文字を含めてください。", field))
}

func (vc *ValidationContext) ValidateContainsSpecialRegx(value, field, errMsg string) {
	re := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`)
	if !re.MatchString(value) {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、特殊文字を含めてください。", field))
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
	if hasNumber {
		return
	}
	if errMsg != "" {
		vc.AddError(field, errMsg)
		return
	}
	vc.AddError(field, fmt.Sprintf("%sには、数字を含めてください。", field))
}

func (vc *ValidationContext) ValidateContainsNumberRegx(value, field, errMsg string) {
	re := regexp.MustCompile(`[0-9]`)
	if !re.MatchString(value) {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、数字を含めてください。", field))
	}
}

// ValidateContainsUppercase checks if the value contains at least one uppercase letter.
func (vc *ValidationContext) ValidateContainsUppercase(value, field, errMsg string) {
	re := regexp.MustCompile(`[A-Z]`)
	if !re.MatchString(value) {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、大文字の英字を含めてください。", field))
	}
}

// ValidateContainsLowercase checks if the value contains at least one lowercase letter.
func (vc *ValidationContext) ValidateContainsLowercase(value, field, errMsg string) {
	re := regexp.MustCompile(`[a-z]`)
	if !re.MatchString(value) {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、小文字の英字を含めてください。", field))
	}
}

// ValidateURL checks if the value is a valid URL.
func (vc *ValidationContext) ValidateURL(value, field, errMsg string) {
	re := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	if !re.MatchString(value) {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効なURLを指定してください。", field))
	}
}

// ValidateFile checks if the value is a valid file path.
func (vc *ValidationContext) ValidateFile(value, field, errMsg string) {
	if _, err := os.Stat(value); err != nil {
		if os.IsNotExist(err) {
			if errMsg != "" {
				vc.AddError(field, errMsg)
				return
			}
			vc.AddError(field, fmt.Sprintf("%sには、有効なファイルパスを指定してください。", field))
		}
	}
}

// ValidateUUID checks if the value is a valid UUID.
func (vc *ValidationContext) ValidateUUID(value, field, errMsg string) {
	if _, err := uuid.Parse(value); err != nil {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効なUUIDを指定してください。", field))
	}
}
