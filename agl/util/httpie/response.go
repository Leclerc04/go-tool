package httpie

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/leclecr04/go-tool/agl/base/trace"
	"github.com/leclecr04/go-tool/agl/util/errs"
)

// Response provides access to a http response.
type Response struct {
	r   *http.Response
	err error
	ctx context.Context
}

// NewResponse wraps aournd http response to make using it easier.
func NewResponse(ctx context.Context, r *http.Response, err error) *Response {
	return &Response{ctx: ctx, r: r, err: errs.Wrap(err)}
}

// Consume consumes the response.
func (r *Response) Consume(action func(r *http.Response) error) error {
	if r.r != nil {
		defer func() {
			// Ensure body is exhausted for connection reuse.
			// We probably don't want to do this if the response left is huge.
			// So use this with good judgement.
			_, err := io.Copy(ioutil.Discard, r.r.Body)
			if err != nil {
				trace.Warningf(r.ctx, "failed to drain body: %v", err)
			}
			err = r.r.Body.Close()
			if err != nil {
				// It is OK to ignore error here, since it doesn't hurt if we
				// cannot close.
				trace.Warningf(r.ctx, "failed to close body: %v", err)
			}
		}()
	}
	if r.err != nil {
		trace.Printf(r.ctx, "http error: %v", r.err)
		return r.err
	}

	if r.r.StatusCode >= 400 {
		b, err := ioutil.ReadAll(r.r.Body)
		if err != nil {
			return errs.Newf("http error: %d, %s, (read buf err: %v) , %s", r.r.StatusCode, r.r.Status, err, string(b))
		}
		return errs.Newf("http error: %d, %s, %s", r.r.StatusCode, r.r.Status, string(b))
	}
	defer trace.Printf(r.ctx, "http done.")
	return action(r.r)
}

// Bytes returns the body as bytes.
func (r *Response) Bytes() ([]byte, error) {
	var b []byte
	err := r.Consume(func(resp *http.Response) error {
		var err error
		b, err = ioutil.ReadAll(resp.Body)
		return err
	})
	return b, err
}

// String returns body as string.
func (r *Response) String() (string, error) {
	b, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// JSON populates JSON.
func (r *Response) JSON(result interface{}) error {
	b, err := r.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, result)
}

func (r *Response) XML(result interface{}) error {
	b, err := r.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, result)
}

func (r *Response) GetStatusCode() (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.r.StatusCode, nil
}

func (r *Response) Header() (http.Header, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.r.Header, nil
}
