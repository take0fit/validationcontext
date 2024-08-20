# ValidationContext Library
ValidationContext is a Go library designed to provide a centralized and efficient approach to managing validations in your applications, especially within a Domain-Driven Design (DDD) context. The library is intended to be easy to use, allowing you to collect validation results from various value objects and handle them collectively at a single point in your application.

## Key Features
- Centralized Validation: Collect and manage validation errors from multiple value objects across your application.
- Deferred Error Handling: Aggregate all validation errors and handle them collectively when necessary.
- Comprehensive Validation Methods: A wide range of built-in validation methods, including file handling, string formatting, and value range checks.
- Customizable: Easily extend the library with custom validation logic to fit specific requirements.
- Enhanced Error Handling: Each validation error now includes a stack trace, and aggregate errors can be retrieved with both error messages and stack traces.

## Use Case
In a DDD context, validations are often dispersed across multiple value objects. ValidationContext allows these validations to be aggregated and handled together, ensuring that all potential issues are addressed before proceeding with business logic.

## Getting Started
### Installation
To install the ValidationContext library, you can use go get:
```sh
go get github.com/yourusername/validationcontext

```
### Example Usage: Validating Multiple Value Objects
Consider a scenario where you need to validate several value objects like HairdresserLicenseImage and UserAddress. ValidationContext helps you collect and handle validation results from these objects at once.
```go
package main

import (
	"fmt"
	"os"

	"github.com/yourusername/validationcontext"
)

type LicenseImage struct {
	File *os.File
}

type UserAddress struct {
	City   string
	Street string
}

func NewLicenseImage(file *os.File, vc *validationcontext.ValidationContext) LicenseImage {
	vc.Required(file, "LicenseImage", " license image is required", false)
	vc.ValidateFileExtension(file, "LicenseImage", []string{".png", ".jpg"}, "Invalid file extension")
	vc.ValidateFileSize(file, "LicenseImage", 2*1024*1024, "File size must be 2MB or less")
	return LicenseImage{File: file}
}

func NewUserAddress(street, city string, vc *validationcontext.ValidationContext) UserAddress {
	vc.Required(street, "Street", "Street is required", false)
	vc.Required(city, "City", "City is required", false)
	return UserAddress{Street: street, City: city}
}

func main() {
	vc := validationcontext.NewValidationContext()

	// Validate LicenseImage
	file, _ := os.Open("test.jpg")
	defer file.Close()
	image := NewLicenseImage(file, vc)

	// Validate UserAddress
	address := NewUserAddress("Main St", "New York", vc)

	// Check if any validation errors occurred
	if vc.HasErrors() {
		// Aggregate and handle all errors at once
		err := vc.AggregateError()
		if err != nil {
			// Print the error messages and stack traces
			fmt.Println(err.Error())
			aggregateErr, ok := err.(*validationcontext.ValidationAggregateError)
			if ok {
				fmt.Println("Stack Traces:")
				for _, trace := range aggregateErr.GetStackTraces() {
					fmt.Println(trace)
				}
			}
		}
	} else {
		fmt.Println("Validation passed for:", image.File.Name(), "and", address)
	}
}

```

## Explanation
- Centralized Validation: The ValidationContext instance (vc) is passed around to each value object, collecting validation errors.
- Deferred Error Handling: After all validations, vc.HasErrors() checks for any errors. If errors exist, they are aggregated and formatted for handling.
- Enhanced Error Information: The updated version includes stack traces for each validation error, which can be accessed via the ValidationAggregateError type

## Validation Methods

ValidationContext provides a variety of built-in validation methods. Below is a table summarizing the available methods:

Method	Description	Example Usage
| Method                      | Description                                                     | Example Usage                                                           |
|-----------------------------|-----------------------------------------------------------------|-------------------------------------------------------------------------|
| Required                    | Ensures a value is not empty or nil                             | `vc.Required(value, "FieldName", "Field is required", false)`           |
| ValidateMinLength           | Checks if a string has at least a certain number of characters  | `vc.ValidateMinLength(value, "FieldName", 5, "Minimum length is 5")`    |
| ValidateMaxLength           | Checks if a string does not exceed a certain number of characters | `vc.ValidateMaxLength(value, "FieldName", 10, "Maximum length is 10")`  |
| ValidateEmail               | Validates if a string is in a proper email format               | `vc.ValidateEmail(email, "Email", "Invalid email format")`              |
| ValidateContainsSpecial     | Ensures a string contains at least one special character        | `vc.ValidateContainsSpecial(value, "FieldName", "Must contain a special character")` |
| ValidateContainsNumber      | Ensures a string contains at least one numeric character        | `vc.ValidateContainsNumber(value, "FieldName", "Must contain a number")`|
| ValidateContainsUppercase   | Ensures a string contains at least one uppercase letter         | `vc.ValidateContainsUppercase(value, "FieldName", "Must contain an uppercase letter")` |
| ValidateContainsLowercase   | Ensures a string contains at least one lowercase letter         | `vc.ValidateContainsLowercase(value, "FieldName", "Must contain a lowercase letter")` |
| ValidateURL                 | Checks if a string is a valid URL                               | `vc.ValidateURL(value, "FieldName", "Invalid URL format")`              |
| ValidateFilePath            | Ensures the file path is valid                                  | `vc.ValidateFilePath(value, "FilePath", "Invalid file path")`           |
| ValidateFileExtension       | Checks if a file has a valid extension                          | `vc.ValidateFileExtension(file, "FieldName", []string{".jpg", ".png"}, "")` |
| ValidateFileSize            | Ensures the file size is within the specified limit             | `vc.ValidateFileSize(file, "FieldName", 2*1024*1024, "File size must be 2MB or less")` |
| ValidateUUID                | Checks if a string is a valid UUID                              | `vc.ValidateUUID(value, "FieldName", "Invalid UUID format")`            |
| ValidateMinValue            | Ensures a numeric value meets the minimum requirement           | `vc.ValidateMinValue(value, "FieldName", 1, "Value must be at least 1")`|
| ValidateMaxValue            | Ensures a numeric value does not exceed the maximum limit       | `vc.ValidateMaxValue(value, "FieldName", 100, "Value must be 100 or less")` |
| ValidateDate                | Ensures a string is a valid date in the format "2006-01-02"     | `vc.ValidateDate(value, "FieldName", "Invalid date format")`            |
| ValidateYearMonth           | Ensures a string is a valid year and month in the format "2006-01" | `vc.ValidateYearMonth(value, "FieldName", "Invalid year-month format")`|
| ValidateYear                | Ensures a string is a valid year                                | `vc.ValidateYear(value, "FieldName", "Invalid year format")`            |
| ValidateMonth               | Ensures a string is a valid month                               | `vc.ValidateMonth(value, "FieldName", "Invalid month format")`          |
| ValidateDateTime            | Ensures a string is a valid date and time in the format "2006-01-02 15:04:05" | `vc.ValidateDateTime(value, "FieldName", "Invalid datetime format")` |
| ValidateTime                | Ensures a string is a valid time in the format "15:04"          | `vc.ValidateTime(value, "FieldName", "Invalid time format")`            |

## Customizing Validation Logic
ValidationContext is designed to be easily extendable, allowing you to implement custom validation logic that fits your specific needs. This can include additional string checks, complex object validations, or even integrating with external validation libraries.

## Conclusion
ValidationContext simplifies the process of managing validations across multiple value objects in your application. It allows you to collect all validation errors in a centralized context and handle them at your convenience, ensuring consistent and comprehensive error management.

This library is ideal for projects that require robust validation mechanisms, particularly in Domain-Driven Design (DDD) contexts where validation logic is scattered across multiple components.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
