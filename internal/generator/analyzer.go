package generator

import (
	"fmt"
	"go/ast"
	"strings"
)

// shouldIncludeMethod determines if a method should be included in the registry
func shouldIncludeMethod(methodName string, targetMethods []string) bool {
	if len(targetMethods) == 0 {
		return strings.HasPrefix(methodName, "New")
	}

	for _, target := range targetMethods {
		if target == methodName {
			return true
		}
	}
	return false
}

// isValidationConstructor checks if a function is a validation constructor
// A validation constructor must have at least one ValidationContext parameter
// AND must not be an entity constructor (with multiple parameters)
func isValidationConstructor(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil || len(fn.Type.Params.List) < 1 {
		return false
	}

	params := fn.Type.Params.List

	// Check if any parameter is ValidationContext
	hasValidationContext := false
	for _, param := range params {
		if isValidationContextType(param.Type) {
			hasValidationContext = true
			break
		}
	}

	if !hasValidationContext {
		return false
	}

	// Exclude entity constructors (functions with more than 2 parameters)
	// Value Object constructors should have exactly 2 parameters: (value, ValidationContext)
	if len(params) > 2 {
		return false
	}

	return true
}

// isValidationContextType checks if an expression represents ValidationContext type
func isValidationContextType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.StarExpr:
		// *validationcontext.ValidationContext case
		if sel, ok := t.X.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				return ident.Name == "validationcontext" && sel.Sel.Name == "ValidationContext"
			}
		}
		// *ValidationContext case (within same package)
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name == "ValidationContext"
		}
	case *ast.SelectorExpr:
		// validationcontext.ValidationContext case (not pointer)
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name == "validationcontext" && t.Sel.Name == "ValidationContext"
		}
	case *ast.Ident:
		// ValidationContext case (within same package, not pointer)
		return t.Name == "ValidationContext"
	}
	return false
}

// extractFirstParamType extracts the type of the first parameter
func extractFirstParamType(fn *ast.FuncDecl) string {
	if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
		return "string"
	}

	firstParam := fn.Type.Params.List[0]
	return typeToString(firstParam.Type)
}

// typeToString converts an AST type expression to string representation
func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "interface{}"
	}
}

// collectRegistrationsFromFile collects constructor registrations from a parsed Go file
func collectRegistrationsFromFile(file *ast.File, config *GenerateConfig, packagePath string, verbose bool) []Registration {
	var registrations []Registration

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if shouldIncludeMethod(fn.Name.Name, config.Methods) && isValidationConstructor(fn) {
				if verbose {
					fmt.Printf("    Including method: %s\n", fn.Name.Name)
				}

				paramType := extractFirstParamType(fn)
				registrationKey := generateRegistrationKey(packagePath, fn.Name.Name)

				reg := Registration{
					ConstructorName: fn.Name.Name,
					RegistrationKey: registrationKey,
					QualifiedName:   fn.Name.Name,
					ParamType:       paramType,
				}

				registrations = append(registrations, reg)
			} else if verbose && shouldIncludeMethod(fn.Name.Name, config.Methods) {
				fmt.Printf("    Skipping method %s (not a validation constructor or entity constructor)\n", fn.Name.Name)
			}
		}
	}

	return registrations
}
