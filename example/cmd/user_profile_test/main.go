package main

import (
	"fmt"

	"github.com/take0fit/validationcontext/example/dto/user/profile"
)

func main() {
	fmt.Println("=== User Profile Validation Test ===\n")

	// Test 1: 正常ケース
	testNormalUserProfile()

	// Test 2: 型変換テスト（Age: int32 → int）
	testTypeConversion()

	// Test 3: バリデーションエラーテスト
	testValidationErrors()
}

func testNormalUserProfile() {
	fmt.Println("--- Test 1: Normal User Profile ---")

	req := &profile.RequestCreateUser{
		FirstName:      "田中",
		LastName:       "太郎",
		AdminFirstName: "管理者",
		AdminLastName:  "花子",
		Age:            int32(25), // int32 → int 変換テスト
	}

	dto, err := profile.NewInputCreateUserDTO(req)
	if err != nil {
		fmt.Printf("✗ Validation error: %v\n", err)
		return
	}

	fmt.Printf("✓ User FirstName: %s\n", dto.UserFirstName.String())
	fmt.Printf("✓ User LastName: %s\n", dto.UserLastName.String())
	fmt.Printf("✓ Admin FirstName: %s\n", dto.AdminFirstName.String())
	fmt.Printf("✓ Admin LastName: %s\n", dto.AdminLastName.String())
	fmt.Printf("✓ Age: %d (converted from int32)\n", dto.Age.Value())
	fmt.Printf("✓ Validation successful!\n\n")
}

func testTypeConversion() {
	fmt.Println("--- Test 2: Type Conversion (int32 → int) ---")

	// 大きなint32値での変換テスト
	req := &profile.RequestCreateUser{
		FirstName:      "大きな",
		LastName:       "年齢",
		AdminFirstName: "管理者",
		AdminLastName:  "テスト",
		Age:            int32(2147483647), // max int32
	}

	dto, err := profile.NewInputCreateUserDTO(req)
	if err != nil {
		fmt.Printf("✗ Type conversion error: %v\n", err)
		return
	}

	fmt.Printf("✓ Large int32 %d successfully converted to int: %d\n", req.Age, dto.Age.Value())
	fmt.Println()
}

func testValidationErrors() {
	fmt.Println("--- Test 3: Validation Errors ---")

	req := &profile.RequestCreateUser{
		FirstName:      "", // 空の名前
		LastName:       "太郎",
		AdminFirstName: "管理者",
		AdminLastName:  "",       // 空の管理者名
		Age:            int32(0), // 無効な年齢
	}

	dto, err := profile.NewInputCreateUserDTO(req)
	if err != nil {
		fmt.Printf("✓ Expected validation errors caught:\n%v\n", err)
	} else {
		fmt.Printf("✗ Should have failed validation, but got: %+v\n", dto)
	}
	fmt.Println()
}
