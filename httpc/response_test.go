package httpc

import (
	"github.com/bellingham07/go-tool/errorx"
	"testing"
)

func TestResponse(t *testing.T) {
	err := errorx.New("test", int(errorx.CodeInternalErr), "test")

	RespError(nil, nil, err)
}
