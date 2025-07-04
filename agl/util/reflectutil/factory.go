package reflectutil

import "reflect"

// MakeFactory returns a factory function that creates new instance of the same type.
// It always returns a pointer to the instance.
func MakeFactory(instance interface{}) func() interface{} {
	typ := reflect.TypeOf(instance)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return func() interface{} {
		return reflect.New(typ).Interface()
	}
}
