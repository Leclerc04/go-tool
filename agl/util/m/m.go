package m

import (
	"sort"

	"reflect"

	"github.com/leclerc04/go-tool/agl/util/jsonutil"
	"github.com/leclerc04/go-tool/agl/util/reflectutil"
)

type M map[string]interface{}
type L []interface{}

func (m M) SortedKeys() []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// StripEmpty returns a new M by removing elements with zero value.
func (m M) StripEmpty() M {
	newM := M{}
	for k, v := range m {
		rv := reflect.ValueOf(v)
		if !reflectutil.IsZero(rv) {
			newM[k] = v
		}
	}
	return newM
}

func DeepCopyM(from interface{}, to interface{}) {
	s := jsonutil.MarshalJSONOrDie(from)
	jsonutil.UnmarshalJSONOrDie([]byte(s), to)
}

func IsM(m interface{}) (M, bool) {
	if m, ok := m.(map[string]interface{}); ok {
		return m, true
	}
	if m, ok := m.(M); ok {
		return m, true
	}
	return nil, false
}
