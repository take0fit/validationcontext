package profile

import "github.com/take0fit/validationcontext"

type Age struct {
	value int
}

//go:generate voauto-gen
func NewAge(v int, vc *validationcontext.ValidationContext) Age {
	vc.Required(v, "Age", "age is required", false)
	return Age{value: v}
}

func (a *Age) Value() int {
	return a.value
}
