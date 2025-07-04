package pipe

import (
	"io"

	"github.com/leclecr04/go-tool/agl/base/concurrent"
	"github.com/leclecr04/go-tool/agl/util/errs"
)

// WriterToReader pipes writer to reader.
func WriterToReader(f func(w io.Writer) error) io.Reader {
	pr, pw := io.Pipe()
	concurrent.GoLite(func() {
		var err error
		defer func() {
			errs.Ignore(pw.CloseWithError(err))
		}()
		err = f(pw)
	})

	return pr
}
