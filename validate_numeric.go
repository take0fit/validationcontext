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
