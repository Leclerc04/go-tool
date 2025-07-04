package ptrutil

import "reflect"

// Elem de-reference given p.
func Elem(p interface{}) interface{} {
	t := reflect.TypeOf(p)       // get p's type
	if t.Kind() != reflect.Ptr { // if the type is NOT point, return p itself
		return p
	}
	v := reflect.ValueOf(p) // otherwise, return p's value
	return v.Elem().Interface()
}

// BoolOrFalse converts a bool pointer to bool.
func BoolOrFalse(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

// StringOrEmpty converts a string pointer to string.
func StringOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// FloatOrZero converts a float pointer to float.
func FloatOrZero(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0.0
}

// FloatOrZero converts a int pointer to int.
func IntOrZero(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

// FloatOrZero converts a int64 pointer to int64.
func Int64OrZero(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}
