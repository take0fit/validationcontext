package main

import (
	"fmt"

	"github.com/take0fit/validationcontext/example/dto/user/profile"
)

func main() {
	req := &profile.RequestCreateUser{
		FirstName:      "田中",
		LastName:       "太郎",
		AdminFirstName: "管理者",
		AdminLastName:  "花子",
	}

	dto, err := profile.NewInputCreateUserDTO(req)
	if err != nil {
		fmt.Println("validation error:", err)
		return
	}

	fmt.Println("User FirstName =", dto.UserFirstName.String())
	fmt.Println("User LastName  =", dto.UserLastName.String())
	fmt.Println("Admin FirstName =", dto.AdminFirstName.String())
	fmt.Println("Admin LastName  =", dto.AdminLastName.String())
	fmt.Println("Validation successful!")
}
