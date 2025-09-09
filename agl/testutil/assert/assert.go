package assert

/*
Custom assertions based on testify/assert.

The name of this package has to be "assert", so that
the testify/assert report will correctly report the caller.
*/

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/leclerc04/go-tool/agl/util/jsonutil"
	"github.com/leclerc04/go-tool/agl/util/must"
	"github.com/leclerc04/go-tool/agl/util/simplejson"
	"github.com/stretchr/testify/assert"
)

// Equal calls assert.Equal.
func Equal(t *testing.T, a, b interface{}) {
	if as, okA := a.(string); okA {
		if bs, okB := b.(string); okB {
			StrEqual(t, as, bs)
			return
		}
	}
	assert.Equal(t, a, b)
}

// NotEqual is assert.NotEqual.
var NotEqual = assert.NotEqual

// NoError is assert.NoError.
var NoError = assert.NoError

// NotEmpty is assert.NotEmpty.
var NotEmpty = assert.NotEmpty

// Contains is assert.Contains.
var Contains = assert.Contains

// True is assert.True.
var True = assert.True

// ErrorHasMessage asserts err.Error() contains a msg as substring.
func ErrorHasMessage(t *testing.T, err error, msg string) {
	if err == nil {
		assert.Fail(t, "Error expected but got nil")
	}
	assert.Contains(t, err.Error(), msg)
}

// StrEqual is a special equal for string.
func StrEqual(t *testing.T, a, b string) {
	d := diff.Diff(a, b)
	if d != "" {
		assert.Fail(t, fmt.Sprintf("Mismatched with diff: \n%s\nActual: \n%s", d, b))
	}
}

// JSONEqual asserts two objects are equal when converted to JSON.
func JSONEqual(t *testing.T, actual interface{}, expected interface{}) {
	s := jsonutil.MarshalJSONOrDie(jsonutil.CleanJSON(expected))
	JSONText(t, jsonutil.CleanJSON(actual), s)
}

// JSONOpt is options for comparing JSON.
type JSONOpt struct {
	IgnoreArrayOrder []string
}

func normalizeJSON(a *simplejson.JSON, opts []JSONOpt) {
	for _, opt := range opts {
		if len(opt.IgnoreArrayOrder) > 0 {
			maybeArr := a.GetPath(opt.IgnoreArrayOrder)
			if maybeArr == nil {
				continue
			}
			arr, ok := maybeArr.CheckArray()
			if !ok {
				continue
			}
			var newArrStr []string
			for _, v := range arr {
				b, err := json.Marshal(v)
				must.Must(err)
				newArrStr = append(newArrStr, string(b))
			}
			sort.Strings(newArrStr)
			var newArr []interface{}
			for _, v := range newArrStr {
				var jv interface{}
				err := json.Unmarshal([]byte(v), &jv)
				must.Must(err)
				newArr = append(newArr, jv)
			}
			a.SetPath(opt.IgnoreArrayOrder, newArr)
		}
	}
}

// JSONText asserts the object has given JSON representation.
func JSONText(t *testing.T, actual interface{}, expected string, opts ...JSONOpt) {
	if expected != "" {
		expectedM := simplejson.New()
		err := json.Unmarshal([]byte(expected), &expectedM)
		assert.NoError(t, err)
		normalizeJSON(expectedM, opts)
		expectedD, err := json.MarshalIndent(expectedM, "", "  ")
		assert.NoError(t, err)
		expected = string(expectedD)
	}

	if actualRaw, isRaw := actual.(jsonutil.RawJSON); isRaw {
		actual = ([]byte)(actualRaw)
	}
	_, ok := actual.([]byte)

	var err error
	if !ok {
		actual, err = json.Marshal(actual)
		assert.NoError(t, err)
	}

	actualB := actual.([]byte)
	actualM := simplejson.New()
	err = json.Unmarshal(actualB, &actualM)
	assert.NoError(t, err)
	normalizeJSON(actualM, opts)

	actualD, err := json.MarshalIndent(actualM, "", "  ")
	assert.NoError(t, err)

	Equal(t, expected, string(actualD))
}

// JSONGet extract values at path, if nil, failed.
func JSONGet(t *testing.T, j *simplejson.JSON, path []string) *simplejson.JSON {
	t.Helper()
	result := j.GetPath(path)
	if result == nil {
		t.Errorf("unexpected json with null at path:\nPath: %v\nJSON: %s", path, jsonutil.MarshalForDebug(j))
	}
	return result
}

// JSONHasStr asserts the object when converted to json, contains given string. Most useful for error assertion.
func JSONHasStr(t *testing.T, actual interface{}, needle string) {
	t.Helper()
	actualB, err := json.Marshal(actual)
	NoError(t, err)
	if needle == "" {
		t.Errorf("needle is empty, actual: %s", jsonutil.MarshalForDebug(actual))
	}
	if !strings.Contains(string(actualB), needle) {
		t.Errorf("given object doesn't contains %s, actual, %s", needle, jsonutil.MarshalForDebug(actual))
	}
}
