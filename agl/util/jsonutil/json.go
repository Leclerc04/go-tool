package jsonutil

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/leclerc04/go-tool/agl/util/errs"
	"github.com/leclerc04/go-tool/agl/util/must"
	"github.com/leclerc04/go-tool/agl/util/simplejson"
)

// RawJSON is the same as json.RawMessage,
// but fix https://github.com/golang/go/issues/14493.
// Do not use json.RawMessage.
type RawJSON []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m RawJSON) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("null"), nil
	}
	return []byte(m), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawJSON) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errs.Internal.Newf("jsonutil.RawJSON: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

var _ encoding.TextUnmarshaler = (*JSONTime)(nil)

func MarshalJSONOrDie(v interface{}) string {
	b, err := json.Marshal(v)
	must.Must(err)
	return string(b)
}

// MarshalOrDie is json.Marshal or die.
func MarshalOrDie(v interface{}) []byte {
	b, err := json.Marshal(v)
	must.Must(err)
	return b
}

func UnmarshalJSONOrDie(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(fmt.Errorf("umarshal failed: %v\n%s", err, string(data)))
	}
}

func CleanIndent(v interface{}) string {
	b, err := json.MarshalIndent(CleanJSON(v), "", "  ")
	must.Must(err)
	return string(b)
}

func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return errs.Internal.Wrapf(err, "source: %s", string(data))
	}
	return nil
}

func CleanJSON(i interface{}) CleanJSONType {
	return CleanJSONType{I: i}
}

// CleanJSONType wrapps an object that is marshaled into json
// where any fields with empty values are removed.
type CleanJSONType struct {
	I interface{}
}

func (j CleanJSONType) IsNil() bool {
	return j.I == nil
}

func (j CleanJSONType) MarshalJSON() ([]byte, error) {
	b0, err := json.Marshal(j.I)
	if err != nil {
		return nil, err
	}
	return StripJSONBytes(b0), nil
}

func StripJSONBytes(b []byte) []byte {
	var r interface{}

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&r)
	if err != nil {
		panic(err)
	}
	newR, _ := jsonStripInterface(r)
	newB, err := json.Marshal(newR)
	if err != nil {
		panic(err)
	}
	return newB
}

func jsonStripInterface(i interface{}) (interface{}, bool) {
	if i == nil {
		return i, true
	}
	switch v := i.(type) {
	case string:
		return v, v == ""
	case []interface{}:
		var newi []interface{}
		for _, ii := range v {
			newii, strip := jsonStripInterface(ii)
			if !strip {
				newi = append(newi, newii)
			}
		}
		return newi, len(newi) == 0
	case json.Number:
		f, err := v.Float64()
		must.Must(err)
		return v, f == 0
	case bool:
		return v, !v
	case map[string]interface{}:
		for key, value := range v {
			newValue, strip := jsonStripInterface(value)
			if strip {
				delete(v, key)
			} else {
				v[key] = newValue
			}
		}
		return v, len(v) == 0
	}
	panic(i)
}

func LogJSON(j interface{}) {
	b, err := json.Marshal(j)
	glog.Infof("%v err: %v", string(b), err)
}

func JSONMergeSample(sample interface{}, target interface{}) interface{} {
	if sample == nil {
		return target
	}
	switch s := sample.(type) {
	case string:
		if _, ok := target.(string); !ok {
			return s
		}
		return target
	case float64:
		if _, ok := target.(float64); !ok {
			return s
		}
		return target
	case bool:
		if _, ok := target.(bool); !ok {
			return s
		}
		return target
	case int:
		if _, ok := target.(int); !ok {
			return s
		}
		return target
	case []interface{}:
		if target, ok := target.([]interface{}); ok {
			for index := range target {
				if len(s) == 0 {
					glog.Infof("JSONMergeSample Warning: Insufficient Sample Information")
				} else {
					target[index] = JSONMergeSample(s[0], target[index])
				}
			}
			return target
		}
		return s
	case map[string]interface{}:
		if target, ok := target.(map[string]interface{}); ok {
			for key, value := range s {
				if _, ok := target[key]; !ok {
					target[key] = value
				} else {
					target[key] = JSONMergeSample(value, target[key])
				}
			}
			return target
		}
		return s
	}
	panic(sample)
}

func ConvertToJSON(data interface{}) (*simplejson.JSON, error) {
	j := simplejson.New()
	err := Remarshal(data, j)
	return j, err
}
