package jsonutil_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	assert2 "github.com/leclerc04/go-tool/agl/testutil/assert"
	. "github.com/leclerc04/go-tool/agl/util/jsonutil"
	"github.com/leclerc04/go-tool/agl/util/must"
)

func TestJSONMergeSample(t *testing.T) {
	// Test case #1 - array of map
	{
		data := []interface{}{
			map[string]interface{}{"number": 1},
			map[string]interface{}{"number": 2},
			map[string]interface{}{"number": 3},
			// a natural case for merging
			// but an unusual case for Populator with sample, since data in db is supposed to be valid all the time.
			// occurs when sample is outdated
			map[string]interface{}{"hello": 2},
		}

		sample := []interface{}{
			map[string]interface{}{"number": 0},
		}

		result := JSONMergeSample(sample, data)

		assert2.JSONText(t, result, `[
				{
					"number": 1
				},
				{
					"number": 2
				},
				{
					"number": 3
				},
				{
					"hello": 2,
					"number": 0
				}
			]`)
	}

	// Test case #2 - some struct
	{
		data := map[string]interface{}{
			"name": "data",
			"letters": []interface{}{
				map[string]interface{}{"name": "aa"},
				map[string]interface{}{"name": "bb"},
			},
			"map_one": map[string]interface{}{"snow": "snowy"},
			"map_two": map[string]interface{}{"sun": "sunny", "rain": "rainy"},
		}

		sample := map[string]interface{}{
			"name": "",
			"letters": []interface{}{
				map[string]interface{}{"name": ""},
			},
			"map_one": "",
			"map_two": map[string]interface{}{"sun": ""},
			"phone":   "",
		}

		result := JSONMergeSample(sample, data)

		assert2.JSONText(t, result, `{
				"letters": [
					{
						"name": "aa"
					},
					{
						"name": "bb"
					}
				],
				"map_one": "",
				"map_two": {
					"rain": "rainy",
					"sun": "sunny"
				},
				"name": "data",
				"phone": ""
			}`)
	}
}

func TestJSONTimeZ(t *testing.T) {
	z := JSONTimeZ{}
	must.Must(z.UnmarshalJSON([]byte(`"2016-07-29T12:34:56Z"`)))
	b, err := z.MarshalJSON()
	if string(b) != `"2016-07-29T12:34:56Z"` {
		t.Error(string(b), err)
	}
}

func TestJSONTimeMicros(t *testing.T) {
	tests := []struct {
		from, to string
	}{
		{`"2016-07-29T12:34:56.987654"`, `"2016-07-29T12:34:56.987654"`},
		{`"2016-07-29T12:34:56"`, `"2016-07-29T12:34:56.000000"`},
	}

	for _, s := range tests {
		z := JSONTimeMicros{}
		err := z.UnmarshalJSON([]byte(s.from))
		assert.NoError(t, err)
		b, err := z.MarshalJSON()
		if string(b) != s.to {
			t.Error(string(b), err)
		}
	}
}

func TestJSONTime(t *testing.T) {
	tests := []struct {
		from, to string
	}{
		{`"2016-07-29T12:34:56.987654"`, `"2016-07-29T12:34:56"`},
		{`"2016-07-29T12:34:56"`, `"2016-07-29T12:34:56"`},
	}

	for _, s := range tests {
		z := JSONTime{}
		err := z.UnmarshalJSON([]byte(s.from))
		assert.NoError(t, err)
		b, err := z.MarshalJSON()
		if string(b) != s.to {
			t.Error(string(b), err)
		}
	}
}

func TestRawJSON(t *testing.T) {
	b := []byte(`{"a": "hello", "k": null}`)
	d := RawJSON{}
	assert.NoError(t, json.Unmarshal(b, &d))

	r, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"hello","k":null}`, string(r))

	r, err = json.Marshal(CleanJSON(d))
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"hello"}`, string(r))
}
