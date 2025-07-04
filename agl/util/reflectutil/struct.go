package reflectutil

import (
	"reflect"
)

// FieldPtrByName given pst *T returns &(pst.field).
// Returns nil if not found.
func FieldPtrByName(pst interface{}, field string) interface{} {
	pstT := reflect.TypeOf(pst)
	if pstT.Kind() != reflect.Ptr || pstT.Elem().Kind() != reflect.Struct {
		panic("given st is not a pointer to a struct")
	}

	pstV := reflect.ValueOf(pst)
	if pstV.IsNil() {
		return nil
	}
	return reflect.Indirect(pstV).FieldByName(field).Addr().Interface()
}
