package errs_test

import (
	"errors"
	"testing"

	"github.com/leclerc04/go-tool/agl/util/errs"
	"github.com/stretchr/testify/assert"
)

var errReal = errors.New("real")

func makeError() error {
	return errs.InvalidArgument.Wrapf(nil, "hi")
}

func TestErrors(t *testing.T) {
	err := makeError()
	assert.Error(t, err)
	assert.Equal(t, "hi", errs.ErrorMessage(err))

	err = errs.Wrap(err)
	assert.Error(t, err)
	assert.True(t, errs.InvalidArgument.Is(err))
	assert.Equal(t, "hi", errs.ErrorMessage(err))

	err = errs.NotFound.Newf("not found")
	assert.True(t, errs.NotFound.Is(err))
	assert.False(t, errs.Forbidden.Is(err))
	assert.Equal(t, "not found", errs.ErrorMessage(err))

	err = errs.Wrap(errReal)
	assert.True(t, err != errReal)
	assert.True(t, errs.Unwrap(err) == errReal)
	assert.Equal(t, "real", errs.ErrorMessage(err))

	err2 := errs.Wrapf(errs.Wrap(errReal), "hello %d", 1)
	assert.True(t, err2 != errReal)
	assert.True(t, err != err2)
	assert.True(t, errs.Unwrap(err2) == errReal)
	assert.Equal(t, "hello 1", errs.ErrorMessage(err2))

	err = errs.InvalidArgument.Wrap(err)
	assert.True(t, errs.Unwrap(err) == errReal)
	assert.True(t, errs.InvalidArgument.Is(err))

	err = errs.InvalidArgument.Newf("cannot update")
	err = errs.WithSubKind(err, "errs_test:test")
	errV, ok := err.(*errs.Error)
	assert.True(t, ok)
	assert.Equal(t, "errs_test:test", errV.SubKind)
	assert.Equal(t, "cannot update", errs.ErrorMessage(err))
}

func TestAttachment(t *testing.T) {
	err := errs.InvalidArgument.WrapcSkipFrame(
		0, nil, "invalid user", "user", 123, "name", "tom",
		"config", map[string]string{"hi": "hi"})
	if err == nil {
		t.Error("err is nil!")
		return
	}
	assert.Equal(t, `InvalidArgument error: invalid user {"config":"map[hi:hi]","name":"tom","user":"123"}.`,
		err.Error())
}
