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
| ValidateIPAddress           | Checks if a value is a valid IP address (IPv4 or IPv6) | `vc.ValidateIPAddress(value, "FieldName", "Invalid IP address")` |
| ValidateIPv4                | Checks if a value is a valid IPv4 address        | `vc.ValidateIPv4(value, "FieldName", "Invalid IPv4 address")` |
| ValidateIPv6                | Checks if a value is a valid IPv6 address        | `vc.ValidateIPv6(value, "FieldName", "Invalid IPv6 address")` |
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

## Automatic Value-Object Binding (`voauto`)

`voauto` is an optional helper package that lets you **declare** how a DTO field
is converted into a Value Object and validated—without writing boiler-plate in
your service layer.

### Installation

Before using `voauto`, you need to install the code generator tool:

#### Option 1: Install from GitHub (Recommended)

    go install github.com/take0fit/validationcontext/cmd/voauto-gen@latest

#### Option 2: Build from Source

    git clone https://github.com/take0fit/validationcontext.git
    cd validationcontext
    go install ./cmd/voauto-gen

Verify the installation:

    voauto-gen -h

### Quick start

#### 1. **Create Value Objects with Generation Comments**

Add `//go:generate voauto-gen` to your constructor functions:

    package profile

    import "github.com/take0fit/validationcontext"

    type FirstName struct {
        value string
    }

    //go:generate voauto-gen
    func NewFirstName(v string, vc *validationcontext.ValidationContext) FirstName {
        vc.Required(v, "FirstName", "first name is required", false)
        return FirstName{value: v}
    }

    func (f *FirstName) String() string { return f.value }

    package profile

    import "github.com/take0fit/validationcontext"

    type LastName struct {
        value string
    }

    //go:generate voauto-gen
    func NewLastName(v string, vc *validationcontext.ValidationContext) LastName {
        vc.Required(v, "LastName", "last name is required", false)
        return LastName{value: v}
    }

    func (l *LastName) String() string { return l.value }

#### 2. **Generate Constructor Registries**

Run the code generator to automatically create registration files:

    # Generate for current directory and subdirectories
    go generate ./...

    # Or run directly for specific directory
    voauto-gen -dir ./domain

    # Verbose output to see what's being processed
    voauto-gen -v -dir ./domain

This will create `registry_init.go` files like:

    // Code generated by voauto-gen; DO NOT EDIT.
    package profile

    import (
        "github.com/take0fit/validationcontext"
        "github.com/take0fit/validationcontext/voauto"
    )

    func init() {
        // Fully qualified key (unique across different packages)
        voauto.Register("user_valueobject_profile_NewFirstName",
            func(v any, vc *validationcontext.ValidationContext) any {
                return NewFirstName(v.(string), vc)
            })
        
        voauto.Register("user_valueobject_profile_NewLastName",
            func(v any, vc *validationcontext.ValidationContext) any {
                return NewLastName(v.(string), vc)
            })
    }

#### 3. **Annotate DTO fields with the `vctag` struct tag**

Format: `vctag:"<ConstructorKey>[,<SourceFieldName>]"`.

    package dto

    import (
        ap "github.com/yourproject/domain/admin/valueobject/profile"
        up "github.com/yourproject/domain/user/valueobject/profile"
        "github.com/take0fit/validationcontext/voauto"
    )

    type RequestCreateUser struct {
        FirstName      string
        LastName       string
        AdminFirstName string
        AdminLastName  string
    }

    type InputCreateUserDTO struct {
        // Manual key specification (fully qualified)
        UserFirstName  up.FirstName `vctag:"user_valueobject_profile_NewFirstName,FirstName"`
        
        // Auto-inference (recommended - automatically determines the correct key)
        UserLastName   up.LastName  `vctag:"auto:NewLastName,LastName"`
        
        // Multiple packages with same constructor names are handled automatically
        AdminFirstName ap.FirstName `vctag:"admin_valueobject_profile_NewFirstName,AdminFirstName"`
        AdminLastName  ap.LastName  `vctag:"auto:NewLastName,AdminLastName"`
    }

    func NewInputCreateUserDTO(req *RequestCreateUser) (*InputCreateUserDTO, error) {
        return voauto.BindAndValidate[InputCreateUserDTO](req)
    }

#### 4. **Bind & validate** in one call:

    func main() {
        req := &RequestCreateUser{
            FirstName:      "田中",
            LastName:       "太郎",
            AdminFirstName: "管理者",
            AdminLastName:  "花子",
        }

        dto, err := NewInputCreateUserDTO(req)
        if err != nil {
            // err can be *validationcontext.ValidationAggregateError
            fmt.Println("Validation error:", err)
            return
        }

        fmt.Println("User FirstName =", dto.UserFirstName.String())
        fmt.Println("User LastName =", dto.UserLastName.String())
        fmt.Println("Admin FirstName =", dto.AdminFirstName.String())
        fmt.Println("Admin LastName =", dto.AdminLastName.String())
        fmt.Println("Validation successful!")
    }

### Generation Command Options

The `voauto-gen` tool supports several command-line options:

    # Basic usage - scan current directory
    voauto-gen

    # Scan specific directory
    voauto-gen -dir ./domain

    # Verbose output (shows processing details)
    voauto-gen -v

    # Custom options in go:generate comments
    //go:generate voauto-gen -output=custom_registry.go
    //go:generate voauto-gen -methods=NewUser,NewProduct
    //go:generate voauto-gen -package=registry

### Tag Syntax Options

#### Auto-inference (Recommended)

    type DTO struct {
        Field ValueType `vctag:"auto:ConstructorName,SourceField"`
    }

The `auto:` prefix automatically infers the full constructor key from the field's type and package path.

#### Manual Key Specification

    type DTO struct {
        Field ValueType `vctag:"full_constructor_key,SourceField"`
    }

Use the exact key as generated in the registry file.

### Multiple Package Support

The generator automatically handles multiple packages with the same constructor names:

    // user/valueobject/profile/name.go
    //go:generate voauto-gen
    func NewFirstName(v string, vc *ValidationContext) FirstName { ... }

    // admin/valueobject/profile/name.go  
    //go:generate voauto-gen
    func NewFirstName(v string, vc *ValidationContext) FirstName { ... }

Both will be registered with unique keys:
- `user_valueobject_profile_NewFirstName`
- `admin_valueobject_profile_NewFirstName`

### Workflow Integration

#### With Go Modules

Add to your `go.mod`:

    module github.com/yourproject

    go 1.24.3

    require (
        github.com/take0fit/validationcontext v1.0.0
    )

#### Build Process Integration

    # In your build script or CI/CD pipeline
    go generate ./...
    go build ./...
    go test ./...

#### IDE Integration

Most Go IDEs support `go:generate` comments. You can:

1. Right-click on the file and select "Generate"
2. Use IDE shortcuts to run `go generate`
3. Configure automatic generation on file save

### Benefits

* **Eliminates boilerplate**: Removes repetitive mapping code from application/service layers.
* **Guarantees validation**: Ensures that **all** validations are executed before business logic runs.
* **Type safety**: Compile-time verification of constructor signatures and types.
* **Safe registration**: Prevents double registration mistakes with automatic key generation.
* **Package isolation**: Supports multiple packages with same constructor names.
* **Auto-inference**: Automatically determines correct constructor keys from field types.

### Troubleshooting

#### Common Issues

1. **"Constructor not found"**:
   - Ensure `//go:generate voauto-gen` is present above your constructor
   - Run `go generate ./...` to regenerate registry files

2. **"Package import errors"**:
   - Check that import paths match your project structure
   - Verify `go.mod` is correctly configured

3. **"Type conversion errors"**:
   - Verify parameter types match between DTO field and constructor
   - Check that constructor signature follows the expected pattern

4. **"Missing registrations"**:
   - Run `voauto-gen -v` for verbose output to debug
   - Check that generated `registry_init.go` files are imported

#### Debug Mode

Use verbose mode to see detailed processing information:

    voauto-gen -v -dir ./domain

This shows:
- Files being processed
- Methods being included/excluded
- Registration keys being generated
- Output files being created

### Best Practices

1. **Consistent structure**: Organize value objects in clear package hierarchies
2. **Use auto-inference**: Prefer `auto:ConstructorName` tags for cleaner code
3. **Run generation regularly**: Include `go generate ./...` in your build process
4. **Version control**: Commit generated `registry_init.go` files
5. **Don't edit generated files**: They will be overwritten on next generation

See `validationcontext/voauto` package documentation for more details.
```

## License
This project is licensed under the MIT License. See the LICENSE file for details.