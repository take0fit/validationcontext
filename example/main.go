package main

import (
	"fmt"

	"github.com/take0fit/validationcontext/example/dto/user/profile"
)

func main() {
	req := &profile.RequestCreateUser{
		FirstName: "FirstName",
		LastName:  "LastName",
	}

	dto, err := profile.NewInputCreateUserDTO(req)
	if err != nil {
		fmt.Println("validation error:", err)
		return
	}

	fmt.Println("FirstName =", dto.FirstName.String())
	fmt.Println("LastName  =", dto.LastName.String())
}
