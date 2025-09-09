package errs

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/leclerc04/go-tool/agl/util/must"
)

// Error implements error interface and capture the relevant information.
type Error struct {
	Kind        Kind
	SubKind     string
	Message     string
	Err         error // Underlying error.
	Attachments map[string]interface{}

	stack stack
}

// Error prints message first, then the top down satck trace.
func (e *Error) Error() string {
	return e.ToString(true)
}

// ToString converts error to a human readable string.
func (e *Error) ToString(relativeStackOnly bool) string {
	return e.toStringInternal(true, relativeStackOnly)
}

// ErrorMessage returns the raw error message without kind or stack.
// If the error is wraped, returns the outer most layer with message set.
func (e *Error) ErrorMessage() string {
	if e.Message != "" {
		if len(e.Attachments) == 0 {
			return e.Message
		}
		return e.Message + " [" + SimpleParamKVToString(e.Attachments) + "]"
	}
	if e.Err != nil {
		if e2, ok := e.Err.(*Error); ok {
			return e2.ToStringNoStack()
		}
		s := e.Err.Error()
		if s != "" {
			return s
		}
	}
	return e.Kind.String() + " error"
}

func (e *Error) ToStringNoStack() string {
	return e.ErrorMessage()
}

func (e *Error) toStringInternal(withStack bool, relativeStackOnly bool) string {
	b := new(bytes.Buffer)

	must.Write(fmt.Fprintf(b, "%s error", e.Kind.String()))
	if e.SubKind != "" {
		must.Write(fmt.Fprintf(b, ": %s subkind", e.SubKind))
	}
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	}
	if len(e.Attachments) > 0 {
		attachStr := ParamKVToString(e.Attachments)
		b.WriteString(" ")
		b.WriteString(attachStr)
	}
	b.WriteString(".")
	if e.Err != nil {
		pad(b, "\n")
		must.Write(fmt.Fprintf(b, "Caused by: %v", e.Err))
	}
	printStack(e.stack, b, relativeStackOnly)

	return b.String()
}

// Wrap wraps the error.
func Wrap(err error) error {
	return WrapfSkipFrame(1, err, "")
}

// Wrapc wraps the error with extra context attachments.
func Wrapc(err error, msg string, kv ...interface{}) error {
	return WrapcSkipFrame(1, err, msg, kv...)
}

func WrapcSkipFrame(depth int, err error, msg string, kv ...interface{}) error {
	if err == nil {
		return nil
	}
	var kind = Internal
	if e2, ok := err.(*Error); ok {
		kind = e2.Kind
	}
	return kind.WrapcSkipFrame(1+depth, err, msg, kv...)
}

// Wrapf is almost equivalent to Internal.Wrapf, but it tries to
// preserve the kind, and ignores nil error.
func Wrapf(err error, msg string, args ...interface{}) error {
	return WrapfSkipFrame(1, err, msg, args...)
}

// WrapfSkipFrame is almost equivalent to Internal.Wrapf, but it tries to
// preserve the kind, and ignores nil error.
func WrapfSkipFrame(depth int, err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if e2, ok := err.(*Error); ok {
		return e2.Kind.WrapfSkipFrame(1+depth, err, msg, args...)
	}
	return Internal.WrapfSkipFrame(1+depth, err, msg, args...)
}

// Newf is equivalent to Internal.Newf.
func Newf(msg string, args ...interface{}) error {
	return Internal.WrapfSkipFrame(1, nil, msg, args...)
}

func New(msg string, kv ...interface{}) error {
	return Internal.WrapcSkipFrame(1, nil, msg, kv...)
}

func WrapPanicValue(v interface{}) error {
	if e, ok := v.(error); ok {
		return Wrap(e)
	}
	return Newf("recover panic from non-error value: %v", v)
}

// WithSubKind 返回一个添加了subKind属性的error,
// 如果传入的err是Error类型则err本身的subKind会被修改并返回原来的err
func WithSubKind(err error, subKind string) error {
	if err == nil {
		return nil
	}
	if !subKindPattern.MatchString(subKind) {
		panic(fmt.Errorf("invalid subKind %s", subKind))
	}
	errV, ok := err.(*Error)
	if !ok {
		errV = Wrap(err).(*Error)
	}
	errV.SubKind = subKind
	return errV
}

// GetSubKind get error's SubKind
func GetSubKind(err error) string {
	if err == nil {
		return ""
	}
	errV, ok := err.(*Error)
	if !ok {
		return ""
	}
	return errV.SubKind
}

// Separator is the string used to separate nested errors. By
// default, to make errors easier on the eye, nested errors are
// indented on a new line. A server may instead choose to keep each
// error on a single line by modifying the separator string, perhaps
// to ":: ".
var Separator = ":\n\t"

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

var subKindPattern = regexp.MustCompile(`^[a-z_]+:[a-zA-Z_0-9]+$`)

// ToStringWithFullStack shows error with full stack.
func ToStringWithFullStack(err error) string {
	if err == nil {
		return ""
	}
	err2, ok := err.(*Error)
	if !ok {
		err2 = Wrap(err).(*Error)
	}
	return err2.ToString(false)
}

// ErrorMessage returns the raw error message without kind or stack.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*Error); ok {
		return e.ErrorMessage()
	}
	return err.Error()
}
