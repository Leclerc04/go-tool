package ptrutil

// Bool creates a pointer of bool.
func Bool(b bool) *bool {
	return &b
}

// Int creates a pointer of int.
func Int(i int) *int {
	return &i
}

// String creates a pointer of string.
func String(v string) *string {
	return &v
}
