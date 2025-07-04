package strs

import (
	"fmt"
	"sort"
	"strings"
)

// StripEmpty returns string slice with empty string stripped.
func StripEmpty(array []string) []string {
	var ret []string
	for _, v := range array {
		if v == "" {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

// RemoveFromSlice returns slices with the given value removed. The order is preserved.
func RemoveFromSlice(value string, array []string) []string {
	var deleted int
	for i := range array {
		j := i - deleted
		if array[j] == value {
			deleted++
			array[j] = array[len(array)-1]
			array = array[:len(array)-1]
		}
	}
	return array
}

// SplitAndTrim splits string by the giving denominator,
// it trims each each parts and removes empty string.
func SplitAndTrim(s string, d string) []string {
	return Filter(
		Map(strings.Split(s, d), strings.TrimSpace),
		func(s string) bool {
			return s != ""
		})
}

// FmtMsgAndArgs returns formated string according to the fmt package.
func FmtMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		if v, ok := msgAndArgs[0].(string); ok {
			return v
		}
		if v, ok := msgAndArgs[0].(error); ok {
			return v.Error()
		}
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

// SetToSlice returns sorted []string in increasing order
func SetToSlice(m map[string]bool) []string {
	result := SetToSliceNoSort(m)
	sort.Strings(result)
	return result
}

// SetToSliceNoSort retrns a string slice from the giving string set.
func SetToSliceNoSort(m map[string]bool) []string {
	var result []string
	for k, v := range m {
		if v {
			result = append(result, k)
		}
	}
	return result
}

// RemoveFromSet removes giving string from the string set.
func RemoveFromSet(m map[string]bool, l ...string) {
	for _, v := range l {
		delete(m, v)
	}
}

// RemoveDuplicate removes duplicate string, and return a new slice
func RemoveDuplicate(list []string) []string {
	seen := map[string]bool{}
	var l []string
	for _, str := range list {
		if seen[str] {
			continue
		}
		seen[str] = true
		l = append(l, str)
	}
	return l
}

// Truncate returns a string with at most n runes.
func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[0:n]
}

// Truncate returns a string with at most n runes.
func TruncateRune(s string, n int) string {
	ru := []rune(s)
	if len(ru) <= n {
		return string(ru)
	}
	return string(ru[0:n])
}
