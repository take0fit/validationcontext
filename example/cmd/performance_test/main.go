package main

import (
	"fmt"
	"time"

	"github.com/take0fit/validationcontext/example/dto/book"
	"github.com/take0fit/validationcontext/example/dto/user/profile"
)

func main() {
	fmt.Println("⚡ PERFORMANCE TESTS")
	fmt.Println(strings.Repeat("=", 40))

	// Test 1: User Profile performance
	testUserProfilePerformance()

	// Test 2: Book Entity performance
	testBookEntityPerformance()

	// Test 3: Bulk operations
	testBulkOperations()
}

func testUserProfilePerformance() {
	fmt.Println("\n📊 User Profile Performance Test")
	fmt.Println(strings.Repeat("-", 30))

	req := &profile.RequestCreateUser{
		FirstName:      "田中",
		LastName:       "太郎",
		AdminFirstName: "管理者",
		AdminLastName:  "花子",
		Age:            int32(25),
	}

	iterations := 10000
	start := time.Now()

	successCount := 0
	for i := 0; i < iterations; i++ {
		_, err := profile.NewInputCreateUserDTO(req)
		if err == nil {
			successCount++
		}
	}

	duration := time.Since(start)
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Average: %v per operation\n", duration/time.Duration(iterations))
	fmt.Printf("Rate: %.2f ops/sec\n", float64(iterations)/duration.Seconds())
}

func testBookEntityPerformance() {
	fmt.Println("\n📊 Book Entity Performance Test")
	fmt.Println(strings.Repeat("-", 30))

	req := &book.RequestCreateBook{
		ID:          int32(123),
		Title:       "Performance Test Book",
		Price:       2500,
		IsAvailable: "true",
		Rating:      4.5,
		PublishedAt: "2023-01-15",
	}

	iterations := 10000
	start := time.Now()

	successCount := 0
	for i := 0; i < iterations; i++ {
		_, err := book.NewCreateBookDTO(req)
		if err == nil {
			successCount++
		}
	}

	duration := time.Since(start)
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Average: %v per operation\n", duration/time.Duration(iterations))
	fmt.Printf("Rate: %.2f ops/sec\n", float64(iterations)/duration.Seconds())
}

func testBulkOperations() {
	fmt.Println("\n📊 Bulk Operations Test")
	fmt.Println(strings.Repeat("-", 30))

	// Create different request variants
	bookReqs := make([]*book.RequestCreateBook, 1000)
	for i := 0; i < 1000; i++ {
		bookReqs[i] = &book.RequestCreateBook{
			ID:          int32(i + 1),
			Title:       fmt.Sprintf("Book #%d", i+1),
			Price:       100 + i,
			IsAvailable: fmt.Sprintf("%t", i%2 == 0),
			Rating:      float64(1 + (i % 5)),
			PublishedAt: "2023-01-01",
		}
	}

	start := time.Now()
	successCount := 0
	errorCount := 0

	for _, req := range bookReqs {
		_, err := book.NewCreateBookDTO(req)
		if err == nil {
			successCount++
		} else {
			errorCount++
		}
	}

	duration := time.Since(start)
	fmt.Printf("Total requests: %d\n", len(bookReqs))
	fmt.Printf("Success: %d\n", successCount)
	fmt.Printf("Errors: %d\n", errorCount)
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Rate: %.2f ops/sec\n", float64(len(bookReqs))/duration.Seconds())
}

func strings.Repeat(s string, count int) string {
	result := make([]byte, 0, len(s)*count)
	for i := 0; i < count; i++ {
		result = append(result, s...)
	}
	return string(result)
}