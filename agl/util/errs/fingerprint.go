package errs

import "bytes"

// Stack return the stack where the error is generated.
func (e *Error) Stack(relative bool) string {
	buf := bytes.Buffer{}
	printStack(e.stack, &buf, true)
	return buf.String()
}
