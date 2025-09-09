package httpc

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/leclerc04/go-tool/errorx"
)

func NewErrorFromRestyResponse(res *resty.Response) *errorx.Error {
	switch res.StatusCode() {
	case http.StatusConflict:

	}
	return errorx.New(
		http.StatusText(res.StatusCode()),
		res.StatusCode(),
		res.String())

}
