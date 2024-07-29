package validationcontext

import (
	"strconv"
)

// ValidateMinValue checks if the value is at least minValue.
func (vc *ValidationContext) ValidateMinValue(value string, field string, minValue int, errMsg string) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		vc.AddError(field, "Invalid integer value")
		return
	}
	if intValue < minValue {
		vc.AddError(field, errMsg)
	}
}

// ValidateMaxValue checks if the value is at most maxValue.
func (vc *ValidationContext) ValidateMaxValue(value string, field string, maxValue int, errMsg string) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		vc.AddError(field, "Invalid integer value")
		return
	}
	if intValue > maxValue {
		vc.AddError(field, errMsg)
	}
}
