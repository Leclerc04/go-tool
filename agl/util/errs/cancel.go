package errs

import (
	"context"
	"net/http"
	"net/url"
)

func deepUnwrap(err error) error {
	if err == nil {
		return nil
	}
	err = Unwrap(err)
	if ue, ok := err.(*url.Error); ok {
		return deepUnwrap(ue.Err)
	}
	return err
}

func IsCancelled(err error) bool {
	err = deepUnwrap(err)
	return err == context.Canceled || err == http.ErrAbortHandler
}
