package httpc

import (
	"testing"

	"github.com/leclecr04/go-tool/errorx"
)

func TestResponse(t *testing.T) {
	err := errorx.New("test", int(errorx.CodeInternalErr), "test")

	RespError(nil, nil, err)
}
