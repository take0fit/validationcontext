package profile

import (
	"github.com/take0fit/validationcontext"
	"github.com/take0fit/validationcontext/voauto"
)

func init() {
	voauto.Register("NewFirstName",
		func(v any, vc *validationcontext.ValidationContext) any {
			return NewFirstName(v.(string), vc)
		})
	voauto.Register("NewLastName",
		func(v any, vc *validationcontext.ValidationContext) any {
			return NewLastName(v.(string), vc)
		})
}
