package reflectutil

import (
	"fmt"
	"reflect"
)

// SliceBuilder helps build a slice of arbitrary type.
type SliceBuilder struct {
	elemStructType reflect.Type
	dstSlice       reflect.Value
	isElemPointer  bool
}

// NewSliceBuilder creates a new SliceBuilder of given type.
// Requires dst to be *[]S, *[]*S for some non-pointer type S.
func NewSliceBuilder(dst interface{}) *SliceBuilder {
	sb := &SliceBuilder{}

	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("dst is not a pointer: %v", dstType))
	}
	dstValue := reflect.ValueOf(dst)
	if dstValue.IsNil() {
		panic("dst cannot be nil")
	}

	sb.dstSlice = dstValue.Elem()
	if sb.dstSlice.Kind() != reflect.Slice {
		panic(fmt.Errorf("dst is not a pointer to slice: %v", dstType))
	}
	elemType := sb.dstSlice.Type().Elem()
	switch elemType.Kind() {
	case reflect.Ptr:
		sb.elemStructType = elemType.Elem()
		sb.isElemPointer = true
	default:
		sb.elemStructType = elemType
	}
	return sb
}

// NewElemPtr returns a new element pointer. It is not added to the slice.
// Call Append to add it.
func (sb *SliceBuilder) NewElemPtr() interface{} {
	ev := reflect.New(sb.elemStructType)
	if sb.elemStructType.Kind() == reflect.Map {
		ev.Elem().Set(reflect.MakeMap(sb.elemStructType))
	}
	return ev.Interface()
}

// Append appends an element pointer to the slice.
func (sb *SliceBuilder) Append(elem interface{}) {
	ev := reflect.ValueOf(elem)
	if !sb.isElemPointer {
		ev = ev.Elem()
	}
	sb.dstSlice.Set(reflect.Append(sb.dstSlice, ev))
}
