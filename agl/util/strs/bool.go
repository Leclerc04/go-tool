package strs

// InSlice returns true if giving string is in the string slice.
func InSlice(value string, array []string) bool {
	for _, s := range array {
		if s == value {
			return true
		}
	}
	return false
}

// HasElementInCommon returns true if there are common member in two arrays
func HasElementInCommon(array1 []string, array2 []string) bool {
	for _, s1 := range array1 {
		for _, s2 := range array2 {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}

// IsSingletonAndEqual returns true if the slice is singleton and it's value euqal to the giving string
func IsSingletonAndEqual(array []string, s string) bool {
	return len(array) == 1 && array[0] == s
}

// IsEmptySlice returns true if the slice does not contain anyting but empty string.
func IsEmptySlice(array []string) bool {
	for _, s := range array {
		if s != "" {
			return false
		}
	}
	return true
}

// EqualSlices tests equality of two string slices.
// It returns true if both content and order of two slices are equal.
func EqualSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
