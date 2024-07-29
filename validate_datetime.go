package validationcontext

import (
	"fmt"
	"time"
)

// ValidateDate checks if the value is a valid date in the format "2006-01-02".
func (vc *ValidationContext) ValidateDate(value, field, errMsg string) {
	if _, err := time.Parse("2006-01-02", value); err != nil {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効な日付を指定してください。", field))
	}
}

// ValidateYearMonth checks if the value is a valid year and month in the format "2006-01".
func (vc *ValidationContext) ValidateYearMonth(value, field, errMsg string) {
	if _, err := time.Parse("2006-01", value); err != nil {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効な年月を指定してください。", field))
	}
}

// ValidateYear checks if the value is a valid year.
func (vc *ValidationContext) ValidateYear(value, field, errMsg string) {
	if _, err := time.Parse("2006", value); err != nil {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効な年を指定してください。", field))
	}
}

// ValidateMonth checks if the value is a valid month.
func (vc *ValidationContext) ValidateMonth(value, field, errMsg string) {
	if _, err := time.Parse("01", value); err != nil {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効な月を指定してください。", field))
	}
}

// ValidateDateTime checks if the value is a valid date and time in the format "2006-01-02 15:04:05".
func (vc *ValidationContext) ValidateDateTime(value, field, errMsg string) {
	if _, err := time.Parse("2006-01-02 15:04:05", value); err != nil {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効な日時を指定してください。", field))
	}
}

// ValidateTime checks if the value is a valid time in the format "15:04".
func (vc *ValidationContext) ValidateTime(value, field, errMsg string) {
	if _, err := time.Parse("15:04", value); err != nil {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sには、有効な時刻を指定してください。", field))
	}
}
