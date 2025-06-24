package voauto

import (
	"fmt"
	"github.com/take0fit/validationcontext"
	"reflect"
	"strings"
	"sync"
)

type ConstructorFunc func(v any, vc *validationcontext.ValidationContext) any

var (
	constructors = make(map[string]ConstructorFunc)
	mu           sync.RWMutex
)

// Register registers a constructor function with a given key
func Register(key string, constructor ConstructorFunc) {
	mu.Lock()
	defer mu.Unlock()
	constructors[key] = constructor
}

// GetConstructor retrieves a constructor function by key
func GetConstructor(key string) ConstructorFunc {
	mu.RLock()
	defer mu.RUnlock()
	return constructors[key]
}

// BindAndValidate binds source data to a DTO and validates it
func BindAndValidate[T any](source any) (*T, error) {
	vc := validationcontext.NewValidationContext()

	var result T
	resultValue := reflect.ValueOf(&result).Elem()
	resultType := reflect.TypeOf(result)

	sourceValue := reflect.ValueOf(source)
	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}

	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		tag := field.Tag.Get("vctag")

		var constructorKey, sourceFieldName string

		if tag != "" {
			// Parse tag
			constructorKey, sourceFieldName = parseTag(tag, field)

			// For auto inference, generate key from type package path
			if strings.HasPrefix(tag, "auto:") {
				constructorKey = generateKeyFromType(field.Type, constructorKey)
			}
		} else {
			// Convention: field name Val -> constructor NewVal
			constructorKey = "New" + field.Name
			sourceFieldName = field.Name
		}

		// Get value from source
		sourceField := sourceValue.FieldByName(sourceFieldName)
		if !sourceField.IsValid() {
			continue
		}

		// Get constructor and execute
		constructor := GetConstructor(constructorKey)
		if constructor == nil {
			return nil, fmt.Errorf("constructor not found: %s", constructorKey)
		}

		resultObj := constructor(sourceField.Interface(), vc)
		if resultObj != nil {
			resultFieldValue := resultValue.Field(i)
			if resultFieldValue.CanSet() {
				resultFieldValue.Set(reflect.ValueOf(resultObj))
			}
		}
	}

	if vc.HasErrors() {
		return nil, vc.AggregateError()
	}

	return &result, nil
}

func parseTag(tag string, field reflect.StructField) (constructorKey, sourceFieldName string) {
	// Remove "auto:" prefix
	cleanTag := strings.TrimPrefix(tag, "auto:")

	parts := strings.Split(cleanTag, ",")
	constructorKey = strings.TrimSpace(parts[0])

	if len(parts) > 1 {
		sourceFieldName = strings.TrimSpace(parts[1])
	} else {
		sourceFieldName = field.Name
	}

	return constructorKey, sourceFieldName
}

func generateKeyFromType(fieldType reflect.Type, constructorName string) string {
	// Get package path
	pkgPath := fieldType.PkgPath()
	if pkgPath == "" {
		return constructorName
	}

	// Get last 2-3 segments from package path
	parts := strings.Split(pkgPath, "/")
	var keyParts []string

	if len(parts) >= 3 {
		keyParts = parts[len(parts)-3:]
	} else if len(parts) >= 2 {
		keyParts = parts[len(parts)-2:]
	} else {
		keyParts = []string{parts[len(parts)-1]}
	}

	packageKey := strings.Join(keyParts, "_")
	return packageKey + "_" + constructorName
}
