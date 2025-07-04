package httpie

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/leclecr04/go-tool/agl/base/mon"
	"github.com/leclecr04/go-tool/agl/base/trace"
)

// httpie encapsulate HTTP request.

var (
	metricRequestsTotal = mon.NewCounterVec(
		"httpie", "requests_total", "Number of http requests sent.",
		[]string{"host", "method", "code"})
)

// Env provides http client feature.
type Env struct {
	Client *http.Client
}

// NewEnv creates a new http Env.
func NewEnv() *Env {
	proxyURL := os.Getenv("AGL_PROXY_URL")
	proxyRulePattern := regexp.MustCompile(os.Getenv("AGL_PROXY_RULE_PATTERN"))
	return &Env{
		Client: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					if proxyURL != "" && proxyRulePattern.Match([]byte(req.URL.Hostname())) {
						return url.Parse("http://" + proxyURL)
					}
					return http.ProxyFromEnvironment(req)
				},
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 5 * time.Minute,
		},
	}
}

// NewEnvWithClient creates a new http env with given client.
func NewEnvWithClient(hc *http.Client) *Env {
	return &Env{
		Client: hc,
	}
}

// IsTimeout returns true if the error is due to timeout.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	switch err := err.(type) {
	case net.Error:
		return err.Timeout()
	}
	return false
}

// Do delegates http.Do.
func (env *Env) Do(ctx context.Context, req *http.Request) *Response {
	trace.Printf(ctx, "http: %s %s", req.Method, req.URL.String())
	r, err := env.Client.Do(req)
	code := "error"
	if err != nil {
		if IsTimeout(err) {
			code = "timeout"
		}
	} else if r != nil {
		code = strconv.Itoa(r.StatusCode)
	}

	metricRequestsTotal.With(mon.Labels{
		"host":   req.URL.Host,
		"method": req.Method,
		"code":   code,
	}).Inc()
	return NewResponse(ctx, r, err)
}

// Get delegates http.Get.
func (env *Env) Get(ctx context.Context, url string) *Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NewResponse(ctx, nil, err)
	}
	return env.Do(ctx, req)
}

// Post delegates http.Post.
func (env *Env) Post(ctx context.Context, url string, bodyType string, body io.Reader) *Response {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return NewResponse(ctx, nil, err)
	}
	req.Header.Set("Content-Type", bodyType)
	return env.Do(ctx, req)
}

// PostForm posts a form data.
func (env *Env) PostForm(ctx context.Context, urlstr string, params url.Values, headers map[string]string) *Response {
	req, err := http.NewRequest("POST", urlstr, strings.NewReader(params.Encode()))
	if err != nil {
		return NewResponse(ctx, nil, err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return env.Do(ctx, req)
}

// PostJSON posts a json.
func (env *Env) PostJSON(ctx context.Context, urlstr string, body interface{}, headers map[string]string, result interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer(bodyBytes))
	if err != nil {
		trace.Printf(ctx, "err: %v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return env.Do(ctx, req).JSON(result)
}

// Post payEase
func (env *Env) PostPayEase(ctx context.Context, url string, bodyType, encryptKey, merchantID, requestID string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	req.Header.Set("encryptKey", encryptKey)
	req.Header.Set("merchantId", merchantID)
	req.Header.Set("requestId", requestID)
	return env.Client.Do(req)
}

// Put delegates http.Put.
func (env *Env) Put(ctx context.Context, url string, bodyType string, body io.Reader) *Response {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return NewResponse(ctx, nil, err)
	}
	req.Header.Set("Content-Type", bodyType)
	return env.Do(ctx, req)
}
