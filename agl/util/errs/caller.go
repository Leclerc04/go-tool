package errs

import (
	"fmt"
	"runtime"
	"strings"
)

// CallerInfo returns file and line number of the caller at given depth.
// Calling with 1 means reporting your caller's location.
func CallerInfo(depth int) string {
	_, file, line, ok := runtime.Caller(1 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := lastIndexN(file, "/", 2)
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func lastIndexN(s, sep string, n int) int {
	idx := -1
	for n > 0 {
		n--
		idx = strings.LastIndex(s, sep)
		if idx < 0 {
			break
		}
		s = s[:idx]
	}
	return idx
}
