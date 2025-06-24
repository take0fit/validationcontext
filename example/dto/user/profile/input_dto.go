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
}

type InputCreateUserDTO struct {
	// 自動推論を使用（型から自動的にパッケージパスを推論）
	UserFirstName  up.FirstName `vctag:"user_valueobject_profile_NewFirstName,FirstName"`
	UserLastName   up.LastName  `vctag:"auto:NewLastName,LastName"`
	AdminFirstName ap.FirstName `vctag:"admin_valueobject_profile_NewFirstName,AdminFirstName"`
	AdminLastName  ap.LastName  `vctag:"auto:NewLastName,AdminLastName"`
}

func NewInputCreateUserDTO(req *RequestCreateUser) (*InputCreateUserDTO, error) {
	return voauto.BindAndValidate[InputCreateUserDTO](req)
}
