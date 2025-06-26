package main

import (
	"fmt"

	"github.com/take0fit/validationcontext/example/dto/user/profile"
)

func main() {
	fmt.Println("🧑‍💼 USER PROFILE TEST")

	// Test 1: Normal case
	fmt.Println("Test 1: Normal User Profile")
	req := &profile.RequestCreateUser{
		FirstName:      "田中",
		LastName:       "太郎",
		AdminFirstName: "管理者",
		AdminLastName:  "花子",
		Age:            int32(25),
	}

	dto, err := profile.NewInputCreateUserDTO(req)
	if err != nil {
		fmt.Printf("  ✗ Error: %v\n", err)
	} else {
		fmt.Printf("  ✓ User: %s %s (Age: %d)\n",
			dto.UserFirstName.String(), dto.UserLastName.String(), dto.Age.Value())
		fmt.Printf("  ✓ Admin: %s %s\n",
			dto.AdminFirstName.String(), dto.AdminLastName.String())
	}

	// Test 2: Validation errors
	fmt.Println("\nTest 2: Validation Errors")
	invalidReq := &profile.RequestCreateUser{
		FirstName:      "",
		LastName:       "太郎",
		AdminFirstName: "管理者",
		AdminLastName:  "",
		Age:            int32(0),
	}

	_, err = profile.NewInputCreateUserDTO(invalidReq)
	if err != nil {
		fmt.Printf("  ✓ Expected validation errors caught: %v\n", err)
	} else {
		fmt.Printf("  ✗ Should have failed validation\n")
	}
}
