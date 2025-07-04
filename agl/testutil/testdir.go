package testutil

import (
	"path"
	"runtime"
	"testing"
)

// GetCurrentSourceFileDir returns the directory of the source file of the caller.
func GetCurrentSourceFileDir(t *testing.T) string {
	_, fn, _, _ := runtime.Caller(1)
	return path.Dir(fn)
}
