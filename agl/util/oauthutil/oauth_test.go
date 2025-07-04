package oauthutil_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/stretchr/testify/assert"

	"net/url"

	"time"

	"github.com/leclecr04/go-tool/agl/util/errs"
	"github.com/leclecr04/go-tool/agl/util/must"
	"github.com/leclecr04/go-tool/agl/util/oauthutil"
)

type testStore struct {
	consumer map[string]string
	token    map[string]string
}

func (s *testStore) GetTokenSecret(token string) (string, error) {
	return s.token[token], nil
}

func (s *testStore) GetConsumerSecret(token string) (string, error) {
	return s.consumer[token], nil
}

func (s *testStore) CheckUniqueness(nonce, timestamp, token string) (bool, error) {
	// fake, always return true.
	return true, nil
}

func TestRoutes(t *testing.T) {
	store := &testStore{
		consumer: map[string]string{
			"e8592a9d-0c99-4772-8972-e2935120c0cc": "YSSc9rex8L1cP9W3IjjFI7TMVU61PkAZTSw2EVZrrBK",
		},
		token: map[string]string{
			"159219fb-41cd-49c0-9e3f-a30e44b18a14": "6oGpLgDiLQ1vCRHpa6EbJYrkObfgr74lRyB3hJt8RVF",
		},
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			fmt.Fprint(w, `{"result": "ok!"}`)
			return
		}
		if r.URL.Path == "/verify_request" {
			or, err := oauthutil.ParseRequest(r)
			if err != nil {
				http.Error(w, err.Error(), errs.GetKind(err).HTTPStatusCode())
				return
			}
			err = oauthutil.Verify(or, store)
			if err != nil {
				http.Error(w, err.Error(), errs.GetKind(err).HTTPStatusCode())
				return
			}
			fmt.Fprint(w, `{"result": "Hello world!"}`)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	})

	port := ":8888"
	go func() {
		must.Must(http.ListenAndServe(port, h))
		// TODO: graceful shutdown when upgraded to go 1.8
	}()
	time.Sleep(300 * time.Millisecond)

	hc := http.DefaultClient
	getURL := func(path string) string {
		return "http://localhost:8888" + path
	}
	// test the server.
	resp, err := hc.Get(getURL("/ping"))
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp, err = hc.Get(getURL("/notfound"))
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	goodToken := &oauth.Credentials{
		Token:  "159219fb-41cd-49c0-9e3f-a30e44b18a14",
		Secret: store.token["159219fb-41cd-49c0-9e3f-a30e44b18a14"],
	}
	params := url.Values{
		"foo":             []string{"bar"},
		"dup":             []string{"du", "pli", "aaa"},
		"percent!@#$%^&*": []string{"val!@#$%^&*"},
		"中文":              []string{"参数值"},
		"empty":           []string{""},
		"":                []string{"empty_key"},
	}

	testGetPost := func(c *oauth.Client, hc *http.Client, token *oauth.Credentials, path string, params url.Values, code int) {
		resp, err := c.Get(hc, token, getURL(path), params)
		assert.NoError(t, err)
		assert.Equal(t, code, resp.StatusCode)
		resp, err = c.Post(hc, token, getURL(path), params)
		assert.NoError(t, err)
		assert.Equal(t, code, resp.StatusCode)
	}

	test := func(method oauth.SignatureMethod) {
		c := &oauth.Client{
			// consumer
			Credentials: oauth.Credentials{
				Token:  "e8592a9d-0c99-4772-8972-e2935120c0cc",
				Secret: store.consumer["e8592a9d-0c99-4772-8972-e2935120c0cc"],
			},
			SignatureMethod: method,
		}
		// simple request without params.
		testGetPost(c, hc, goodToken, "/verify_request", url.Values{}, 200)

		// with params
		testGetPost(c, hc, goodToken, "/verify_request", params, 200)

		// no token
		testGetPost(c, hc, nil, "/verify_request", params, 200)

		// bad token
		badToken := &oauth.Credentials{
			Token:  "159219fb-41cd-49c0-9e3f-a30e44b18a14",
			Secret: "bad_token",
		}
		testGetPost(c, hc, badToken, "/verify_request", params, 403)
		badToken = &oauth.Credentials{
			Token:  "not found token",
			Secret: "",
		}
		testGetPost(c, hc, badToken, "/verify_request", params, 403)

		// bad consumer
		c.Credentials = oauth.Credentials{
			Token:  "e8592a9d-0c99-4772-8972-e2935120c0cc",
			Secret: "bad_consumer",
		}
		testGetPost(c, hc, goodToken, "/verify_request", params, 403)
		c.Credentials = oauth.Credentials{
			Token:  "not found consumer",
			Secret: "",
		}
		testGetPost(c, hc, goodToken, "/verify_request", params, 403)
		c.Credentials = oauth.Credentials{
			Token:  "",
			Secret: "",
		}
		testGetPost(c, hc, goodToken, "/verify_request", params, 400)
	}

	test(oauth.HMACSHA1)
	test(oauth.PLAINTEXT)
}
