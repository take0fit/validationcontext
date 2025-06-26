package validationcontext

import (
	"fmt"
)

func (vc *ValidationContext) ValidateMinValue(value int, field string, minValue int, errMsg string) {
	if value < minValue {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sは%d以上で入力してください。", field, minValue))
	}
}

func (vc *ValidationContext) ValidateMaxValue(value int, field string, maxValue int, errMsg string) {
	if value > maxValue {
		if errMsg != "" {
			vc.AddError(field, errMsg)
			return
		}
		vc.AddError(field, fmt.Sprintf("%sは%d以下で入力してください。", field, maxValue))
	}
}

// ValidateMinFloatValue validates that a float64 value meets the minimum requirement
func (vc *ValidationContext) ValidateMinFloatValue(value float64, field string, min float64, message string) {
	if value < min {
		vc.AddError(field, message)
	}
}

// ValidateMaxFloatValue validates that a float64 value does not exceed the maximum limit
func (vc *ValidationContext) ValidateMaxFloatValue(value float64, field string, max float64, message string) {
	if value > max {
		vc.AddError(field, message)
	}
}

// ValidateFloatRange validates that a float64 value is within the specified range
func (vc *ValidationContext) ValidateFloatRange(value float64, field string, min, max float64, message string) {
	if value < min || value > max {
		vc.AddError(field, message)
	}
}
