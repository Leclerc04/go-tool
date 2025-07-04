package errs

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func BuildParamKV(kv ...interface{}) map[string]interface{} {
	var attach map[string]interface{}
	var invalidValues []interface{}
	for i := 0; i < len(kv); i += 2 {
		if i+1 >= len(kv) {
			invalidValues = append(invalidValues, kv[i])
			continue
		}

		k, ok := kv[i].(string)
		if !ok {
			invalidValues = append(invalidValues, k)
			invalidValues = append(invalidValues, kv[i+1])
			continue
		}
		if attach == nil {
			attach = map[string]interface{}{}
		}
		attach[k] = kv[i+1]
	}
	if len(invalidValues) > 0 {
		if attach == nil {
			attach = map[string]interface{}{}
		}
		attach["_invalid_attached"] = invalidValues
	}
	{
		for k, v := range attach {
			if valueIsEmptyForPrinting(v) {
				delete(attach, k)
				continue
			}
			attach[k] = fmt.Sprintf("%+v", v)
		}
	}
	return attach
}

func valueIsEmptyForPrinting(v interface{}) bool {
	if v == nil {
		return true
	}
	if vv, ok := v.(string); ok && vv == "" {
		return true
	}
	return false
}

func maybeQuote(v interface{}) string {
	s := fmt.Sprint(v)
	if !strings.Contains(s, " ") {
		return s
	}
	return fmt.Sprintf("%q", s)
}

func SimpleParamKVToString(kv map[string]interface{}) string {
	if len(kv) == 0 {
		return ""
	}
	var ss []string
	var keys []string
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := kv[k]
		ss = append(ss, fmt.Sprintf("%s=%s", k, maybeQuote(v)))
	}
	return strings.Join(ss, " ")
}

func ParamKVToString(kv map[string]interface{}) string {
	if len(kv) == 0 {
		return ""
	}
	ab, err := json.Marshal(kv)
	if err != nil {
		return fmt.Sprintf("%+v", kv)
	}
	return string(ab)
}
