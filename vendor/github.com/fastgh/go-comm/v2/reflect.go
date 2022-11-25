package comm

import "reflect"

// Copied From http://golang.org/src/encoding/json/encode.go
// Lines 280 - 296
func IsEmptyReflectValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
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

func IsEmptyValue(v any) bool {
	return IsEmptyReflectValue(reflect.ValueOf(v))
}

func IsPrimitiveReflectValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsPrimitiveReflectValue(v.Elem())
	default:
		return false
	}
}

func IsPrimitiveValue(v any) bool {
	return IsPrimitiveReflectValue(reflect.ValueOf(v))
}
