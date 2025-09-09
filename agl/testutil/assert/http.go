package assert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/leclerc04/go-tool/agl/util/deepcopy"
	"github.com/leclerc04/go-tool/agl/util/must"
	"github.com/stretchr/testify/assert"
)

// HTTPGet is a convenient method for HTTP. Prefer to use HTTP direclty in new code.
func HTTPGet(t *testing.T, h http.Handler, path string, params map[string]string, result interface{}) {
	HTTP(path).Return(result).QueryParams(params).Do(t, h)
}

// HTTPGetError is a convenient method for HTTP. Prefer to use HTTP direclty in new code.
func HTTPGetError(t *testing.T, h http.Handler, code int, path string, params map[string]string) {
	HTTP(path).QueryParams(params).ExpectCode(code).Do(t, h)
}

// HTTPPost is a convenient method for HTTP. Prefer to use HTTP direclty in new code.
func HTTPPost(t *testing.T, h http.Handler, path string, params map[string]string, body, result interface{}) {
	HTTP(path).Post().JSON(body).Return(result).QueryParams(params).Do(t, h)
}

// HTTPPostError is a convenient method for HTTP. Prefer to use HTTP direclty in new code.
func HTTPPostError(t *testing.T, h http.Handler, code int, path string, params map[string]string, body, result interface{}) {
	HTTP(path).Post().JSON(body).QueryParams(params).ExpectCode(code).Do(t, h)
}

// HTTPParams encode optional params for HTTP method.
type HTTPParams struct {
	Method          string
	Path            string
	Query           map[string]string
	FormValues      map[string]string
	ExpectedCode    int
	Body            bytes.Buffer
	ContentType     string
	Result          interface{}
	CookieJar       *cookiejar.Jar
	Host            string
	FinalPath       *string
	Header          http.Header
	RequestModifier func(*http.Request)
	Ctx             context.Context
}

// Visit sets the path.
func (p HTTPParams) Visit(path string) HTTPParams {
	p.Path = path
	return p
}

// FormDataField specifies a single form data field.
type FormDataField struct {
	FieldValue string
	FileName   string
	File       []byte
}

// FormData attaches form data.
func (p HTTPParams) FormData(fields map[string]FormDataField) HTTPParams {
	func() {
		writer := multipart.NewWriter(&p.Body)
		defer func() {
			must.Must(writer.Close())
		}()
		p.ContentType = writer.FormDataContentType()

		for k, v := range fields {
			if len(v.File) > 0 {
				fileWirter, err := writer.CreateFormFile(k, v.FileName)
				must.Must(err)
				must.Write(fileWirter.Write(v.File))
				continue
			}
			must.Must(writer.WriteField(k, v.FieldValue))
		}
	}()
	if p.Method == "" {
		p.Method = "POST"
	}
	return p
}

// File attaches a file to the request.
func (p HTTPParams) File(fieldName, fileName string, fileContent []byte) HTTPParams {
	return p.FormData(map[string]FormDataField{
		fieldName: FormDataField{
			FileName: fileName,
			File:     fileContent,
		},
	})
}

// JSON attaches a json body to the request.
func (p HTTPParams) JSON(j interface{}) HTTPParams {
	must.Must(json.NewEncoder(&p.Body).Encode(j))
	p.ContentType = "application/json"
	if p.Method == "" {
		p.Method = "POST"
	}
	return p
}

func (p HTTPParams) SetContext(ctx context.Context) HTTPParams {
	p.Ctx = ctx
	return p
}

// SetFormBody attach a post form to the request.
func (p HTTPParams) SetFormBody(v url.Values) HTTPParams {
	p.Body.WriteString(v.Encode())
	p.ContentType = "application/x-www-form-urlencoded"
	if p.Method == "" {
		p.Method = "POST"
	}
	return p
}

// QueryParam sets a single query param.
func (p HTTPParams) QueryParam(key, value string) HTTPParams {
	if p.Query == nil {
		p.Query = map[string]string{}
	}
	p.Query[key] = value
	return p
}

// QueryParams sets multipe query params.
func (p HTTPParams) QueryParams(vs map[string]string) HTTPParams {
	if p.Query == nil {
		p.Query = vs
		return p
	}
	for k, v := range vs {
		p.Query[k] = v
	}
	return p
}

// ModifyRequest sets request modifier.
func (p HTTPParams) ModifyRequest(modifier func(*http.Request)) HTTPParams {
	p.RequestModifier = modifier
	return p
}

// HTTP starts a http request builder.
func HTTP(path string) HTTPParams {
	return HTTPParams{Path: path}
}

// Return sets the result.
func (p HTTPParams) Return(r interface{}) HTTPParams {
	p.Result = r
	return p
}

// WithCookieJar access cookies from the jar.
func (p HTTPParams) WithCookieJar(jar *cookiejar.Jar) HTTPParams {
	p.CookieJar = jar
	return p
}

// SetHost set host in request.
func (p HTTPParams) SetHost(host string) HTTPParams {
	p.Host = host
	return p
}

// Post changes method to POST.
func (p HTTPParams) Post() HTTPParams {
	p.Method = "POST"
	return p
}

// Put changes method to PUT.
func (p HTTPParams) Put() HTTPParams {
	p.Method = "PUT"
	return p
}

// Delete changes method to DELETE.
func (p HTTPParams) Delete() HTTPParams {
	p.Method = "DELETE"
	return p
}

// ExpectCode sets the expected status code.
func (p HTTPParams) ExpectCode(c int) HTTPParams {
	p.ExpectedCode = c
	return p
}

// RedirectTo returns the final path after redirections.
func (p HTTPParams) RedirectTo(s *string) HTTPParams {
	p.FinalPath = s
	return p
}

// SetHeader sets header in request.
func (p HTTPParams) SetHeader(header http.Header) HTTPParams {
	p.Header = deepcopy.Iface(header).(http.Header)
	return p
}

// Do makes a http request and verify some common behavior.
func (p HTTPParams) Do(t *testing.T, h http.Handler) {
	if p.Method == "" {
		p.Method = "GET"
	}

	header := http.Header{}
	if p.Header != nil {
		header = p.Header
	}
	if p.ContentType != "" {
		header.Add("Content-Type", p.ContentType)
	}
	var w *httptest.ResponseRecorder
	for i := 0; i < 10; i++ {
		urlToVisit := buildURL(p.Path, p.Query)
		host := "dummyhost"
		if p.Host != "" {
			host = p.Host
		}
		dummyURL, err := url.Parse("https://" + host + urlToVisit)
		must.Must(err)

		req, err := http.NewRequest(p.Method, urlToVisit, &p.Body)
		if p.Ctx != nil {
			req = req.WithContext(p.Ctx)
		}
		assert.NoError(t, err, "failed creating request")

		if p.Host != "" {
			req.Host = p.Host
		}

		req.Header = header
		if p.CookieJar != nil {
			for _, c := range p.CookieJar.Cookies(dummyURL) {
				req.AddCookie(c)
				// A shortcut: Perform CSRF for the client.
				if c.Name == "csrftoken" {
					v, err := url.QueryUnescape(c.Value)
					must.Must(err)
					req.Header.Set("X-CSRFToken", v)
				}
			}
		}

		w = httptest.NewRecorder()
		if p.RequestModifier != nil {
			p.RequestModifier(req)
		}
		h.ServeHTTP(w, req)

		if w.Code == 301 || w.Code == 302 {
			p.Path = w.Header().Get("Location")
			u, err := url.Parse(p.Path)
			if err != nil {
				assert.NoError(t, err)
				break
			}
			if u.Host != "" {
				if p.FinalPath != nil {
					*p.FinalPath = p.Path
				}
				break
			}
			// log.Printf("[assert] Redirect to %s", path)
			if i == 9 {
				assert.Fail(t, "too many redirect and the testutil is unable to handle")
				break
			}
			continue
		}

		if p.FinalPath != nil {
			*p.FinalPath = p.Path
		}

		if p.ExpectedCode > 0 {
			assert.Equal(t, p.ExpectedCode, w.Code, "status code not matched")
		} else if w.Code >= 400 {
			assert.Fail(t, fmt.Sprintf("http %s %s failed: %v %s", p.Method, urlToVisit, w.Code, w.Body.String()))
		}

		if w.Code < 400 && p.CookieJar != nil {
			req := http.Request{}
			req.Header = http.Header{"Cookie": w.Result().Header["Set-Cookie"]}
			p.CookieJar.SetCookies(dummyURL, req.Cookies())
		}
		break
	}

	if p.Result != nil {
		var err error
		bs := w.Body.Bytes()
		if bytePtr, ok := p.Result.(*[]byte); ok {
			*bytePtr = []byte(bs)
		} else if strPtr, ok := p.Result.(*string); ok {
			*strPtr = string(bs)
		} else {
			err = json.Unmarshal(bs, p.Result)
		}
		if w.Code == p.ExpectedCode || (p.ExpectedCode == 0 && w.Code < 400) {
			assert.NoError(t, err, "unmarshal http result failed, raw: "+string(bs))
		}
	}
}

func buildURL(path string, params map[string]string) string {
	ret := path
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	if len(values) > 0 {
		ret += "?" + values.Encode()
	}
	return ret
}

// NewRequest creates an HTTP request.
func NewRequest(method, path string) *http.Request {
	r, err := http.NewRequest(method, path, nil)
	must.Must(err)
	return r
}

// DoHTTP just perform a plain request and return the response.
func DoHTTP(ctx context.Context, h http.Handler, req *http.Request) *http.Response {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w.Result()
}

// RequestSummary is a summary of http.Request.
type RequestSummary struct {
	Header http.Header
	Method string
	URL    string
}

// NewRequestSummary creates a RequestSummary.
func NewRequestSummary(req *http.Request) RequestSummary {
	return RequestSummary{
		Header: req.Header,
		Method: req.Method,
		URL:    req.URL.String(),
	}
}
