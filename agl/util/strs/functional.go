package strs

// Map returns new slice with each element applied by the given function.
func Map(s []string, f func(string) string) []string {
	if s == nil {
		return nil
	}
	o := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		o[i] = f(s[i])
	}
	return o
}

// Filter returns new slice with each element matched with given function.
func Filter(s []string, f func(string) bool) []string {
	if s == nil {
		return nil
	}
	var o []string
	for _, si := range s {
		if f(si) {
			o = append(o, si)
		}
	}
	return o
}
