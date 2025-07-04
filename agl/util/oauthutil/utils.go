package oauthutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// percentEncode percent encodes a string according to RFC 5849 3.6
func percentEncode(src string) string {
	return string(encode(src, false))
}

func signatureBaseStringPlainText(consumerSecret, tokenSecret string) string {
	return fmt.Sprintf("%s&%s", percentEncode(consumerSecret), percentEncode(tokenSecret))
}

func baseStringURI(host, path string, isHTTPS bool) string {
	var scheme string
	if isHTTPS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	filteredHost := func() string {
		splitted := strings.Split(host, ":")
		switch {
		case len(splitted) == 2 && splitted[1] == "80" && scheme == "http":
			return splitted[0]
		case len(splitted) == 2 && splitted[1] == "443" && scheme == "https":
			return splitted[0]
		default:
			return host
		}
	}()

	ur := &url.URL{
		Scheme: scheme,
		Host:   strings.ToLower(filteredHost),
		Path:   path,
	}

	return ur.String()
}

func formatRequestParametersForSigning(values ...url.Values) string {
	var kvs []kv
	for _, vals := range values {
		for k, v := range vals {
			for _, s := range v {
				kvs = append(kvs, kv{
					key:   percentEncode(k),
					value: percentEncode(s),
				})
			}
		}
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Less(kvs[j])
	})
	params := make([]string, len(kvs))
	for i := range kvs {
		params[i] = fmt.Sprint(kvs[i].key, "=", kvs[i].value)
	}

	return strings.Join(params, "&")
}

// kv is an internal representation of a key value pair, used to sort the list
type kv struct{ key, value string }

func (k kv) Less(that kv) bool {
	// 1.0a/9.1.1 states that kvp must be sorted by key, then by value,
	if k.key == that.key {
		return k.value < that.value
	}
	return k.key < that.key
}

func splitParams(p url.Values) (oauthParams url.Values, otherParams url.Values) {
	oauthParams = make(url.Values)
	otherParams = make(url.Values)

	for k, v := range p {
		if strings.HasPrefix(k, "oauth_") {
			oauthParams[k] = v
		} else {
			otherParams[k] = v
		}
	}
	return
}

func checkHMACSHA1(message, key, signature string) bool {
	hashfun := hmac.New(sha1.New, []byte(key))
	_, err := hashfun.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	sig := base64.StdEncoding.EncodeToString(hashfun.Sum(nil))
	return safeEqual(signature, sig)
}

// safeEqual uses ConstantTimeCompare to prevent against timing attack.
func safeEqual(x string, y string) bool {
	return subtle.ConstantTimeCompare([]byte(x), []byte(y)) == 1
}
