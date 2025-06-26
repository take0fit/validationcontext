package profile

import (
	ap "github.com/take0fit/validationcontext/example/domain/admin/valueobject/profile"
	up "github.com/take0fit/validationcontext/example/domain/user/valueobject/profile"
	"github.com/take0fit/validationcontext/voauto"
)

type RequestCreateUser struct {
	FirstName      string
	LastName       string
	AdminFirstName string
	AdminLastName  string
	Age            int32
}

type InputCreateUserDTO struct {
	UserFirstName  up.FirstName `vctag:"user_valueobject_profile_NewFirstName,FirstName"`
	UserLastName   up.LastName  `vctag:"auto:NewLastName,LastName"`
	AdminFirstName ap.FirstName `vctag:"admin_valueobject_profile_NewFirstName,AdminFirstName"`
	AdminLastName  ap.LastName  `vctag:"auto:NewLastName,AdminLastName"`
	Age            up.Age       `vctag:"auto:NewAge,Age"`
}

func NewInputCreateUserDTO(req *RequestCreateUser) (*InputCreateUserDTO, error) {
	return voauto.BindAndValidate[InputCreateUserDTO](req)
}
