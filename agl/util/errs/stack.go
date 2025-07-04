package errs

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type stack []uintptr

func callers(skip int) stack {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(3+skip, pcs[:])
	var st stack = pcs[0:n]
	return st
}

func mergeStack(e *Error) {
	e2, ok := e.Err.(*Error)
	if !ok {
		return
	}

	// Move distinct callers from inner error to outer error
	// (and throw the common callers away)
	// so that we only print the stack trace once.
	i := 0
	toMerge := false
	for ; i < len(e.stack) && i < len(e2.stack); i++ {
		if e.stack[len(e.stack)-1-i] != e2.stack[len(e2.stack)-1-i] {
			break
		}
		toMerge = true
	}
	if !toMerge {
		return
	}
	if ok { // The stacks have some PCs in common.
		head := e2.stack[:len(e2.stack)-i]
		tail := e.stack
		e.stack = make([]uintptr, len(head)+len(tail))
		copy(e.stack, head)
		copy(e.stack[len(head):], tail)
		e2.stack = nil
	}
}

// GetStack returns the current call stack without some uninteresting frames.
func GetStack(depth int) string {
	b := bytes.Buffer{}
	printStack(callers(depth+1), &b, false)
	return b.String()
}

// printStack prints stack topdown.
func printStack(s stack, b *bytes.Buffer, relativeStackOnly bool) {
	if len(s) == 0 {
		return
	}
	printCallers := callers(0)

	// Iterate backward through e.callers (the last in the stack is the
	// earliest call, such as main) skipping over the PCs that are shared
	// by the error stack and by this function call stack, printing the
	// names of the functions and their file names and line numbers in
	// reverse order.
	var prev string // the name of the last-seen function
	diff := !relativeStackOnly
	var lines []string

	lastFrame := frame(s, len(s)-1)
	frameSkipper := NewFrameSkipper(lastFrame.File)
	skippingFrames := 0
	for i := 0; i < len(s); i++ {
		thisFrame := frame(s, i)
		name := thisFrame.Func.Name()

		if !diff && i < len(printCallers) {
			if name == frame(printCallers, i).Func.Name() {
				// both stacks share this PC, skip it.
				continue
			}
			// No match, don't consider printCallers again.
			diff = true
		}

		// Don't print the same function twice.
		// (Can happen when multiple error stacks have been coalesced.)
		if name == prev {
			continue
		}
		skip := frameSkipper.ShouldSkip(thisFrame.File)
		if skip {
			skippingFrames++
			continue
		}
		if skippingFrames > 0 {
			lines = append(lines, fmt.Sprintf("<%d frames omitted>", skippingFrames))
			skippingFrames = 0
		}

		// Find the uncommon prefix between this and the previous
		// function name, separating by dots and slashes.
		trim := 0
		for {
			j := strings.IndexAny(name[trim:], "./")
			if j < 0 {
				break
			}
			if !strings.HasPrefix(prev, name[:j+trim]) {
				break
			}
			trim += j + 1 // skip over the separator
		}

		// Set to lines
		dots := "..."
		if trim == 0 {
			dots = ""
		}
		lines = append(lines, fmt.Sprintf("%v:%d: %s%s", thisFrame.File, thisFrame.Line, dots, name[trim:]))

		prev = name
	}

	// Do the printing
	for i := len(lines) - 1; i >= 0; i-- {
		pad(b, Separator)
		fmt.Fprint(b, lines[i])
	}
}

// frame returns the nth frame, with the frame at top of stack being 0.
func frame(callers []uintptr, n int) *runtime.Frame {
	frames := runtime.CallersFrames(callers)
	var f runtime.Frame
	for i := len(callers) - 1; i >= n; i-- {
		var ok bool
		f, ok = frames.Next()
		if !ok {
			break // Should never happen, and this is just debugging.
		}
	}
	return &f
}
