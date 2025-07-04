package ints

func IntInSlice(value int, array []int) bool {
	for _, s := range array {
		if s == value {
			return true
		}
	}
	return false
}

func Int64InSlice(value int64, array []int64) bool {
	for _, s := range array {
		if s == value {
			return true
		}
	}
	return false
}
