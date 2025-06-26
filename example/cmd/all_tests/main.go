package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/take0fit/validationcontext/example/dto/book"
	"github.com/take0fit/validationcontext/example/dto/user/profile"
)

func main() {
	fmt.Println("===============================================")
	fmt.Println("  ValidationContext Library - All Tests")
	fmt.Println("===============================================")

	// User Profile Tests
	runUserProfileTests()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Book Entity Tests
	runBookEntityTests()

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// Advanced Type Conversion Tests
	runAdvancedTypeConversionTests()

	fmt.Println("\n===============================================")
	fmt.Println("  All Tests Completed")
	fmt.Println("===============================================")
}

func runUserProfileTests() {
	fmt.Println("🧑‍💼 USER PROFILE TESTS")
	fmt.Println(strings.Repeat("-", 30))

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
		fmt.Printf("  ✓ Expected validation errors caught\n")
	} else {
		fmt.Printf("  ✗ Should have failed validation\n")
	}
}

func runBookEntityTests() {
	fmt.Println("📚 BOOK ENTITY TESTS")
	fmt.Println(strings.Repeat("-", 30))

	// Test 1: Normal case with type conversions
	fmt.Println("Test 1: Normal Book Creation with Type Conversions")
	req := &book.RequestCreateBook{
		ID:          int32(123), // int32 → int
		Title:       "Go Programming",
		Price:       2500,   // int → float64
		IsAvailable: "true", // "true" → bool
		Rating:      4.8,    // float64 → int (4)
		PublishedAt: "2023-01-15",
	}

	dto, err := book.NewCreateBookDTO(req)
	if err != nil {
		fmt.Printf("  ✗ Error: %v\n", err)
	} else {
		fmt.Printf("  ✓ Book ID: %d (int32 → int)\n", dto.ID.Value())
		fmt.Printf("  ✓ Title: %s\n", dto.Title.String())
		fmt.Printf("  ✓ Price: $%.2f (int → float64)\n", dto.Price.Value())
		fmt.Printf("  ✓ Available: %t (\"true\" → bool)\n", dto.IsAvailable.Value())
		fmt.Printf("  ✓ Rating: %d/5 (%.1f → int)\n", dto.Rating.Value(), req.Rating)
		fmt.Printf("  ✓ Published: %s\n", dto.PublishedAt.String())
	}

	// Test 2: JSON simulation
	fmt.Println("\nTest 2: JSON Data Simulation")
	jsonData := `{
		"id": 456,
		"title": "Advanced Go Patterns",
		"price": 3500,
		"is_available": "false",
		"rating": 4.2,
		"published_at": "2023-03-10"
	}`

	var jsonReq book.RequestCreateBook
	if err := json.Unmarshal([]byte(jsonData), &jsonReq); err != nil {
		fmt.Printf("  ✗ JSON parse error: %v\n", err)
		return
	}

	jsonDto, err := book.NewCreateBookDTO(&jsonReq)
	if err != nil {
		fmt.Printf("  ✗ JSON validation error: %v\n", err)
	} else {
		fmt.Printf("  ✓ JSON → DTO conversion successful\n")
		fmt.Printf("  ✓ Book: %s ($%.2f, Rating: %d)\n",
			jsonDto.Title.String(), jsonDto.Price.Value(), jsonDto.Rating.Value())
	}

	// Test 3: Validation errors
	fmt.Println("\nTest 3: Book Validation Errors")
	invalidReq := &book.RequestCreateBook{
		ID:          int32(-1),      // Invalid ID
		Title:       "",             // Empty title
		Price:       -100,           // Negative price
		IsAvailable: "maybe",        // Invalid bool
		Rating:      10.0,           // Invalid rating
		PublishedAt: "invalid-date", // Invalid date
	}

	_, err = book.NewCreateBookDTO(invalidReq)
	if err != nil {
		fmt.Printf("  ✓ Multiple validation errors caught\n")
	} else {
		fmt.Printf("  ✗ Should have failed validation\n")
	}
}

func runAdvancedTypeConversionTests() {
	fmt.Println("🔄 ADVANCED TYPE CONVERSION TESTS")
	fmt.Println(strings.Repeat("-", 30))

	// Test 1: String to Bool variations
	fmt.Println("Test 1: String → Bool Conversions")
	boolTests := []struct {
		input    string
		expected bool
		hasError bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"1", true, false},
		{"0", false, false},
		{"TRUE", true, false},
		{"FALSE", false, false},
		{"invalid", false, true},
	}

	for _, test := range boolTests {
		req := &book.RequestCreateBook{
			ID:          int32(1),
			Title:       "Test Book",
			Price:       100,
			IsAvailable: test.input,
			Rating:      3.0,
			PublishedAt: "2023-01-01",
		}

		dto, err := book.NewCreateBookDTO(req)
		if test.hasError {
			if err != nil {
				fmt.Printf("  ✓ '%s' → Expected error\n", test.input)
			} else {
				fmt.Printf("  ✗ '%s' → Should have failed\n", test.input)
			}
		} else {
			if err != nil {
				fmt.Printf("  ✗ '%s' → Unexpected error: %v\n", test.input, err)
			} else {
				actual := dto.IsAvailable.Value()
				if actual == test.expected {
					fmt.Printf("  ✓ '%s' → %t\n", test.input, actual)
				} else {
					fmt.Printf("  ✗ '%s' → Expected %t, got %t\n", test.input, test.expected, actual)
				}
			}
		}
	}

	// Test 2: Numeric precision tests
	fmt.Println("\nTest 2: Numeric Precision & Boundary Tests")
	numericTests := []struct {
		name   string
		rating float64
		expect int
	}{
		{"Exact integer", 4.0, 4},
		{"Round down", 4.9, 4},
		{"Minimum", 1.1, 1},
		{"Maximum", 4.99, 4},
	}

	for _, test := range numericTests {
		req := &book.RequestCreateBook{
			ID:          int32(1),
			Title:       "Numeric Test",
			Price:       100,
			IsAvailable: "true",
			Rating:      test.rating,
			PublishedAt: "2023-01-01",
		}

		dto, err := book.NewCreateBookDTO(req)
		if err != nil {
			fmt.Printf("  ✗ %s: Error %v\n", test.name, err)
		} else {
			actual := dto.Rating.Value()
			if actual == test.expect {
				fmt.Printf("  ✓ %s: %.1f → %d\n", test.name, test.rating, actual)
			} else {
				fmt.Printf("  ✗ %s: %.1f → Expected %d, got %d\n", test.name, test.rating, test.expect, actual)
			}
		}
	}

	// Test 3: Large number conversions
	fmt.Println("\nTest 3: Large Number Conversions")
	req := &book.RequestCreateBook{
		ID:          int32(2147483647), // max int32
		Title:       "Large Numbers Test",
		Price:       999999, // large int → float64
		IsAvailable: "1",
		Rating:      5.0,
		PublishedAt: "2023-12-31",
	}

	dto, err := book.NewCreateBookDTO(req)
	if err != nil {
		fmt.Printf("  ✗ Large number test failed: %v\n", err)
	} else {
		fmt.Printf("  ✓ Max int32 → int: %d\n", dto.ID.Value())
		fmt.Printf("  ✓ Large int → float64: $%.2f\n", dto.Price.Value())
		fmt.Printf("  ✓ \"1\" → bool: %t\n", dto.IsAvailable.Value())
	}
}
