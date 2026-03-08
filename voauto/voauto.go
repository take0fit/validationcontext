package voauto

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/take0fit/validationcontext"
)

type ConstructorFunc func(v any, vc *validationcontext.ValidationContext) any

var (
	constructors = make(map[string]ConstructorFunc)
	mu           sync.RWMutex
)

// Register registers a constructor function with a given key.
func Register(key string, constructor ConstructorFunc) {
	mu.Lock()
	defer mu.Unlock()
	constructors[key] = constructor
}

// GetConstructor retrieves a constructor function by key.
func GetConstructor(key string) ConstructorFunc {
	mu.RLock()
	defer mu.RUnlock()
	return constructors[key]
}

// BindAndValidate binds source data to a target struct and validates it.
func BindAndValidate[T any](source any) (*T, error) {
	vc := validationcontext.NewValidationContext()

	sourceValue, err := normalizeSourceValue(source)
	if err != nil {
		return nil, err
	}

	var result T
	resultValue := reflect.ValueOf(&result).Elem()
	resultType := reflect.TypeOf(result)
	if resultType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("target type must be struct: %s", resultType.Kind())
	}

	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		constructorKey, sourceFieldName, parseErr := resolveBindingRule(field)
		if parseErr != nil {
			return nil, parseErr
		}

		sourceField := sourceValue.FieldByName(sourceFieldName)
		if !sourceField.IsValid() {
			return nil, fmt.Errorf("source field not found: %s", sourceFieldName)
		}

		constructor := GetConstructor(constructorKey)
		if constructor == nil {
			return nil, fmt.Errorf("constructor not found: %s", constructorKey)
		}

		resultObj, callErr := callConstructor(constructor, sourceField.Interface(), vc, constructorKey)
		if callErr != nil {
			return nil, callErr
		}
		if resultObj == nil {
			continue
		}

		resultFieldValue := resultValue.Field(i)
		if !resultFieldValue.CanSet() {
			return nil, fmt.Errorf("target field cannot be set: %s", field.Name)
		}

		resultObjValue := reflect.ValueOf(resultObj)
		if resultObjValue.Type().AssignableTo(resultFieldValue.Type()) {
			resultFieldValue.Set(resultObjValue)
			continue
		}
		if resultObjValue.Type().ConvertibleTo(resultFieldValue.Type()) {
			resultFieldValue.Set(resultObjValue.Convert(resultFieldValue.Type()))
			continue
		}
		return nil, fmt.Errorf("constructor result type mismatch for field %s: got %s, want %s", field.Name, resultObjValue.Type(), resultFieldValue.Type())
	}

	if vc.HasErrors() {
		return nil, vc.AggregateError()
	}

	return &result, nil
}

func normalizeSourceValue(source any) (reflect.Value, error) {
	if source == nil {
		return reflect.Value{}, fmt.Errorf("source must not be nil")
	}

	sourceValue := reflect.ValueOf(source)
	if sourceValue.Kind() == reflect.Ptr {
		if sourceValue.IsNil() {
			return reflect.Value{}, fmt.Errorf("source pointer must not be nil")
		}
		sourceValue = sourceValue.Elem()
	}
	if sourceValue.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("source must be a struct or pointer to struct: %s", sourceValue.Kind())
	}
	return sourceValue, nil
}

func resolveBindingRule(field reflect.StructField) (constructorKey, sourceFieldName string, err error) {
	tag := strings.TrimSpace(field.Tag.Get("vctag"))
	if tag == "" {
		return "New" + field.Name, field.Name, nil
	}

	constructorKey, sourceFieldName = parseTag(tag, field)
	if constructorKey == "" {
		return "", "", fmt.Errorf("invalid vctag for field %s: constructor key is empty", field.Name)
	}
	if sourceFieldName == "" {
		return "", "", fmt.Errorf("invalid vctag for field %s: source field is empty", field.Name)
	}
	if strings.HasPrefix(tag, "auto:") {
		constructorKey = generateKeyFromType(field.Type, constructorKey)
	}
	return constructorKey, sourceFieldName, nil
}

func callConstructor(constructor ConstructorFunc, input any, vc *validationcontext.ValidationContext, constructorKey string) (result any, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("constructor panic: %s: %v", constructorKey, recovered)
		}
	}()
	result = constructor(input, vc)
	return result, nil
}

func parseTag(tag string, field reflect.StructField) (constructorKey, sourceFieldName string) {
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
	pkgPath := fieldType.PkgPath()
	if pkgPath == "" {
		return constructorName
	}

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
