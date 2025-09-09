package oauthutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"regexp"

	"strconv"

	"strings"

	"github.com/leclerc04/go-tool/agl/util/errs"
	"github.com/leclerc04/go-tool/agl/util/strs"
	"github.com/leclerc04/go-tool/agl/util/timeutil"
)

const (
	oauthSignatureParamKey       = "oauth_signature"
	oauthSignatureMethodParamKey = "oauth_signature_method"
	oauthConsumerKeyParamKey     = "oauth_consumer_key"
	oauthTokenParamKey           = "oauth_token"
	oauthTimestampParamKey       = "oauth_timestamp"
	oauthNonceParamKey           = "oauth_nonce"
	oauthVersionParamKey         = "oauth_version"

	oauthDefaultVersion      = "1.0"
	signatureMethodPlainText = "PLAINTEXT"
	signatureMethodHmacSha1  = "HMAC-SHA1"
	signatureMethodRsaSha1   = "RSA-SHA1"
	authorizationPrefix      = "OAuth " // trailing space is intentional
)

var (
	headerParamPattern = regexp.MustCompile(`(\w+)="(.*?)"`)
	TimestampThreshold = 600 // inherited from django.
)

type Store interface {
	GetConsumerSecret(token string) (string, error)
	GetTokenSecret(token string) (string, error)
	CheckUniqueness(nonce, timestamp, token string) (bool, error)
}

type Request struct {
	HTTPRequest *http.Request
	realm       string

	// the signature is not stored in oauthParameter but in this field for easier signature generation.
	signature string

	oauthParameters url.Values
	otherParameters url.Values
}

// ParseRequest parses an oauth request from Authorization header, content body, or url query.
func ParseRequest(req *http.Request) (*Request, error) {
	if req.Method != "GET" && req.Method != "POST" {
		return nil, errs.InvalidArgument.Newf("Request method not supported: %s", req.Method)
	}
	r := &Request{
		HTTPRequest:     req,
		oauthParameters: make(url.Values),
		otherParameters: make(url.Values),
	}

	oauthHeaderParams, err := r.parseAuthorizationHeader()
	if err != nil {
		return nil, err
	}
	oauthBodyParams, otherBodyParams, err := r.parseRequestBody()
	if err != nil {
		return nil, err
	}
	oauthQueryParams, otherQueryParams := r.parseRequestURLQuery()

	switch {
	case len(oauthHeaderParams) > 0:
		r.oauthParameters = oauthHeaderParams
	case len(oauthBodyParams) > 0:
		r.oauthParameters = oauthBodyParams
	case len(oauthQueryParams) > 0:
		r.oauthParameters = oauthQueryParams
	default:
		return nil, errs.InvalidArgument.Newf("Oauth params not provided through either authorization header, request body, or url query.")
	}

	if len(otherBodyParams) > 0 && len(otherQueryParams) > 0 {
		return nil, errs.InvalidArgument.Newf("Does not support parameter passed through both body and url query.")
	}
	if len(otherBodyParams) > 0 {
		r.otherParameters = otherBodyParams
	}
	if len(otherQueryParams) > 0 {
		r.otherParameters = otherQueryParams
	}

	// remove the oauth signature, and move it to signature field
	r.signature = r.getOAuthParameter(oauthSignatureParamKey)
	r.oauthParameters.Del(oauthSignatureParamKey)
	return r, nil
}

// Verify checks timestamp, nonce and signature.
func Verify(r *Request, store Store) error {
	// check version
	if r.Version() != oauthDefaultVersion {
		return errs.InvalidArgument.Newf("Only oauth 1.0 is supported.")
	}

	// check consumer key
	if r.ConsumerKey() == "" {
		return errs.InvalidArgument.Newf("Consumer key not set.")
	}

	// check if has valid secrets
	consumerSecret, err := store.GetConsumerSecret(r.ConsumerKey())
	if err != nil {
		return err
	}
	if consumerSecret == "" {
		return errs.Forbidden.Newf("consumer secret")
	}
	tokenSecret, err := store.GetTokenSecret(r.Token())
	if err != nil {
		return err
	}
	if r.Token() != "" && tokenSecret == "" {
		return errs.Forbidden.Newf("token secret")
	}

	// check timestamp
	// RFC 5849 section 3.1[1]: timestamp or nonce on PLAINTEXT is optional.
	if !(r.SignatureMethod() == signatureMethodPlainText && r.Timestamp() == "") {
		ts, err := strconv.Atoi(r.Timestamp())
		if err != nil {
			return errs.InvalidArgument.Newf("timestamp invalid: %s Error: %v", r.Timestamp(), err)
		}
		now := int(timeutil.Now().Unix())
		if now-ts > TimestampThreshold {
			return errs.Forbidden.Newf("timestamp")
		}
	}

	// check nonce
	if !(r.SignatureMethod() == signatureMethodPlainText && (r.Timestamp() == "" || r.Nonce() == "")) {
		unique, err := store.CheckUniqueness(r.Nonce(), r.Timestamp(), r.Token())
		if err != nil {
			return err
		}
		if !unique {
			return errs.Forbidden.Newf("nonce")
		}
	}

	// check signature
	method := r.SignatureMethod()
	switch method {
	case signatureMethodPlainText:
		if !safeEqual(
			signatureBaseStringPlainText(consumerSecret, tokenSecret),
			r.signature) {
			return errs.Forbidden.Newf("token secret")
		}
		return nil
	case signatureMethodHmacSha1:
		if !checkHMACSHA1(
			r.signatureBaseString(),
			signatureBaseStringPlainText(consumerSecret, tokenSecret),
			r.signature) {
			return errs.Forbidden.Newf("token secret")
		}
		return nil
	case signatureMethodRsaSha1:
		return errs.InvalidArgument.Newf("RSA-SHA1 signature method not supported.")
	default:
		return errs.InvalidArgument.Newf("Unknown signature method: %s", method)
	}
}

func (r *Request) signatureBaseString() string {
	return fmt.Sprintf("%s&%s&%s",
		strings.ToUpper(r.HTTPRequest.Method),
		percentEncode(baseStringURI(r.HTTPRequest.Host, r.HTTPRequest.URL.Path, r.HTTPRequest.TLS != nil)),
		percentEncode(formatRequestParametersForSigning(r.oauthParameters, r.otherParameters)),
	)
}

// parseAuthorizationHeader parses authorization header.
//
//	For example:
//
//	  Authorization: OAuth realm="Example",
//	     oauth_consumer_key="0685bd9184jfhq22",
//	     oauth_token="ad180jjd733klru7",
//	     oauth_signature_method="HMAC-SHA1",
//	     oauth_signature="wOJIO9A2W5mFwDgiDvZbTSMK%2FPY%3D",
//	     oauth_timestamp="1487711237",
//	     oauth_nonce="4572616e48616d6d65724c61686176",
//	     oauth_version="1.0"
func (r *Request) parseAuthorizationHeader() (url.Values, error) {
	header := r.HTTPRequest.Header.Get("Authorization")
	if !strings.HasPrefix(header, authorizationPrefix) {
		return nil, errs.InvalidArgument.Newf("Invalid Authorization Header: " + header)
	}

	header = strings.TrimPrefix(header, authorizationPrefix)
	params := strs.SplitAndTrim(header, ",")

	values := make(url.Values)
	for _, p := range params {
		pargs := headerParamPattern.FindStringSubmatch(p)
		if len(pargs) < 3 {
			return nil, errs.InvalidArgument.Newf("Invalid Authorization Header: %s", header)
		}
		decoded, err := url.QueryUnescape(pargs[2])
		if err != nil {
			return nil, errs.InvalidArgument.Newf("Header value cannot be decoded. Header: %s Error: %v", header, err)
		}
		key := pargs[1]
		if key == "realm" {
			r.realm = decoded
			continue
		}
		if !strings.HasPrefix(key, "oauth_") {
			return nil, errs.InvalidArgument.Newf("Invalid oauth key: %s, header %s", key, header)
		}
		values.Set(key, decoded)
	}
	return values, nil
}

// parseRequestBody parses request body if content type is application/x-www-form-urlencoded
// returns oauth params and other params.
func (r *Request) parseRequestBody() (oauthParams url.Values, otherParams url.Values, err error) {
	if r.HTTPRequest.Method != "POST" {
		return nil, nil, nil
	}
	ct, _, err := mime.ParseMediaType(r.HTTPRequest.Header.Get("Content-Type"))
	if err != nil || ct != "application/x-www-form-urlencoded" {
		return nil, nil, nil
	}

	b, err := ioutil.ReadAll(r.HTTPRequest.Body)
	if err != nil {
		return nil, nil, errs.InvalidArgument.Newf("Body cannot be read: %v", err)
	}
	defer func() {
		// reinitialize Body with ReadCloser over the []byte
		r.HTTPRequest.Body = ioutil.NopCloser(bytes.NewReader(b))
	}()

	values, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, nil, errs.InvalidArgument.Newf("Body cannot be decoded: %v", err)
	}

	oauthParams, otherParams = splitParams(values)
	return
}

// parseRequestURLQuery parses request url query, returns oauth params and other params.
func (r *Request) parseRequestURLQuery() (oauthParams url.Values, otherParams url.Values) {
	return splitParams(r.HTTPRequest.URL.Query())
}

// getOAuthParameter retrieve any oauth parameter.
func (r *Request) getOAuthParameter(name string) string {
	return r.oauthParameters.Get(name)
}

func (r *Request) SignatureMethod() string {
	return r.getOAuthParameter(oauthSignatureMethodParamKey)
}

func (r *Request) ConsumerKey() string {
	return r.getOAuthParameter(oauthConsumerKeyParamKey)
}

func (r *Request) Token() string {
	return r.getOAuthParameter(oauthTokenParamKey)
}

func (r *Request) Timestamp() string {
	return r.getOAuthParameter(oauthTimestampParamKey)
}

func (r *Request) Nonce() string {
	return r.getOAuthParameter(oauthNonceParamKey)
}

func (r *Request) Version() string {
	return r.getOAuthParameter(oauthVersionParamKey)
}
