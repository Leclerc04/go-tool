package strs

import (
	"bytes"
	"io"
	"io/ioutil"
)

// UntouchingRead allows f to read from a io.Reader, returning the io.Reader as
// if it has not been read.
func UntouchingRead(reader io.Reader, f func(r io.Reader)) io.Reader {
	var buf bytes.Buffer
	f(io.TeeReader(reader, &buf))
	return io.MultiReader(&buf, reader)
}

type untouchingReadCloser struct {
	io.Reader
	rc io.ReadCloser
}

func (r untouchingReadCloser) Close() error {
	return r.rc.Close()
}

// UntouchingReadCloser is the similar to UntouchingRead but also handle close.
func UntouchingReadCloser(reader io.ReadCloser, f func(r io.Reader)) io.ReadCloser {
	return untouchingReadCloser{
		Reader: UntouchingRead(reader, f),
		rc:     reader,
	}
}

// UntouchingReadAll read from reader, and return a reader as if the read didn't happen.
func UntouchingReadAll(reader io.ReadCloser, result *[]byte, err *error) io.ReadCloser {
	if reader == nil {
		return nil
	}
	return UntouchingReadCloser(reader, func(r io.Reader) {
		*result, *err = ioutil.ReadAll(r)
	})
}
