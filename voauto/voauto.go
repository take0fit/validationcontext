package voauto

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/take0fit/validationcontext"
)

// Constructor represents an application-provided factory.
type Constructor func(any, *validationcontext.ValidationContext) any

// registry stores factories found by Bind().
var registry = map[string]Constructor{}

// onceMap guarantees each factory is registered only once.
var (
	onceMap = map[string]*sync.Once{}
	mu      sync.Mutex // protects registry & onceMap
)

// Register adds a factory, but only on the first call for that name.
func Register(name string, c Constructor) {
	mu.Lock()
	o, ok := onceMap[name]
	if !ok {
		o = &sync.Once{}
		onceMap[name] = o
	}
	mu.Unlock()

	o.Do(func() {
		mu.Lock()
		registry[name] = c
		mu.Unlock()
	})
}

// Bind populates dto fields through registered factories.
func Bind(dto any, src any, vc *validationcontext.ValidationContext) error {
	dstVal := reflect.ValueOf(dto)
	if dstVal.Kind() != reflect.Pointer || dstVal.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dto must be pointer to struct")
	}
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Pointer {
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return fmt.Errorf("src must be struct or pointer to struct")
	}

	dstStruct := dstVal.Elem()
	dstType := dstStruct.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)

		tag := field.Tag.Get("vctag")
		var ctorName string
		srcFieldName := field.Name

		if tag != "" {
			parts := strings.Split(tag, ",")
			ctorName = strings.TrimSpace(parts[0])
			if ctorName == "" {
				return fmt.Errorf("vctag of field %s: constructor name is required", field.Name)
			}
			if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
				srcFieldName = strings.TrimSpace(parts[1])
			}
		} else {
			ctorName = "New" + field.Name
		}

		srcField := srcVal.FieldByName(srcFieldName)
		if !srcField.IsValid() {
			continue
		}

		mu.Lock()
		ctor, ok := registry[ctorName]
		mu.Unlock()
		if !ok {
			return fmt.Errorf("constructor %s not registered", ctorName)
		}

		vo := ctor(srcField.Interface(), vc)
		dstStruct.Field(i).Set(reflect.ValueOf(vo))
	}
	return nil
}

// BindAndValidate is a convenience wrapper.
func BindAndValidate[T any](src any) (*T, error) {
	dto := new(T)
	vc := validationcontext.NewValidationContext()

	if err := Bind(dto, src, vc); err != nil {
		return nil, err
	}
	if vc.HasErrors() {
		return nil, vc.AggregateError()
	}
	return dto, nil
}
