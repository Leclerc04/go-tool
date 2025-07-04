package strs

import (
	"strings"

	"github.com/leclecr04/go-tool/agl/util/randutil"
)

const randStringPool = "abcdefghijklmnopqrstuvwxyz0123456789"

// GenerateRandomStringURLSafe returns a URL-safe string of given length.
// The string is lower case only to make it easier for human.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) string {
	return GenerateRandomString(randStringPool, n)
}

// GenerateRandomString returns a function that can generate random strings
// of any length.
func GenerateRandomString(dict string, l int) string {
	var builder strings.Builder
	builder.Grow(l)
	for i := 0; i < l; i++ {
		builder.WriteByte(dict[randutil.Rand.Intn(len(dict))])
	}
	return builder.String()
}
