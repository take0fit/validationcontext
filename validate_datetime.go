package validationcontext

import (
	"fmt"
	"time"
)

func (vc *ValidationContext) validateTimeLayout(value, layout, field, errMsg, defaultMsg string) {
	if _, err := time.Parse(layout, value); err != nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、%s", field, defaultMsg))
	}
}

// ValidateDate checks if the value is a valid date in the format "2006-01-02".
func (vc *ValidationContext) ValidateDate(value, field, errMsg string) {
	vc.validateTimeLayout(value, "2006-01-02", field, errMsg, "有効な日付を指定してください。")
}

// ValidateYearMonth checks if the value is a valid year and month in the format "2006-01".
func (vc *ValidationContext) ValidateYearMonth(value, field, errMsg string) {
	vc.validateTimeLayout(value, "2006-01", field, errMsg, "有効な年月を指定してください。")
}

// ValidateYear checks if the value is a valid year.
func (vc *ValidationContext) ValidateYear(value, field, errMsg string) {
	vc.validateTimeLayout(value, "2006", field, errMsg, "有効な年を指定してください。")
}

// ValidateMonth checks if the value is a valid month.
func (vc *ValidationContext) ValidateMonth(value, field, errMsg string) {
	vc.validateTimeLayout(value, "01", field, errMsg, "有効な月を指定してください。")
}

// ValidateDateTime checks if the value is a valid date and time in the format "2006-01-02 15:04:05".
func (vc *ValidationContext) ValidateDateTime(value, field, errMsg string) {
	vc.validateTimeLayout(value, "2006-01-02 15:04:05", field, errMsg, "有効な日時を指定してください。")
}

// ValidateTime checks if the value is a valid time in the format "15:04".
func (vc *ValidationContext) ValidateTime(value, field, errMsg string) {
	vc.validateTimeLayout(value, "15:04", field, errMsg, "有効な時刻を指定してください。")
}
