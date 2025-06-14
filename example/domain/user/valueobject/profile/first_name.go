package profile

import "github.com/take0fit/validationcontext"

type FirstName struct {
	value string
}

func NewFirstName(v string, vc *validationcontext.ValidationContext) FirstName {
	vc.Required(v, "FirstName", "first name is required", false)
	return FirstName{value: v}
}

func (f *FirstName) String() string { return f.value }
