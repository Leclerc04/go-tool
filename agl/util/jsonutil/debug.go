package jsonutil

import (
	"encoding/json"
	"fmt"
)

// MarshalForDebug marshals an interface{} to string for debug proposes. It might not return
// a valid JSON.
func MarshalForDebug(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", v)
	}
	s := string(j)
	if len(s) > 4096 {
		s = s[:4096]
	}
	return s
}
