package profile

import "github.com/take0fit/validationcontext"

type LastName struct {
	value string
}

//go:generate voauto-gen
func NewLastName(v string, vc *validationcontext.ValidationContext) LastName {
	vc.Required(v, "LastName", "last name is required", false)
	return LastName{value: v}
}

func (l *LastName) String() string { return l.value }
