package validationcontext

import (
	"fmt"
	"reflect"
)

// Required adds a required validation rule to the context.
func (vc *ValidationContext) Required(value interface{}, field string, message string, skipNil bool) {
	value, isNil := indirect(value)
	if skipNil && !isNil && isEmpty(value) || !skipNil && (isNil || isEmpty(value)) {
		if message == "" {
			message = fmt.Sprintf("%sは必須項目です。", field)
		}
		vc.AddError(field, message)
	}
}

// indirect returns the value, after dereferencing as many times as necessary to reach the base type.
func indirect(value interface{}) (interface{}, bool) {
	if value == nil {
		return nil, true
	}
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, true
		}
		v = v.Elem()
	}
	return v.Interface(), false
}

// isEmpty checks if a value is considered empty.
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
