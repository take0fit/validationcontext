package profile

import (
	"github.com/take0fit/validationcontext/example/domain/user/valueobject/profile"
	"github.com/take0fit/validationcontext/voauto"
)

type RequestCreateUser struct {
	FirstName string
	LastName  string
}

type InputCreateUserDTO struct {
	FirstName profile.FirstName `vctag:"NewFirstName,FirstName"` // ctor, src
	LastName  profile.LastName  `vctag:"NewLastName,LastName"`
}

func NewInputCreateUserDTO(req *RequestCreateUser) (*InputCreateUserDTO, error) {
	return voauto.BindAndValidate[InputCreateUserDTO](req)
}
