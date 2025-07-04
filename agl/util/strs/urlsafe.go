package strs

import (
	"github.com/mr-tron/base58/base58"
)

// EncodeURLSafe converts bytes to a string that can be put in URL.
func EncodeURLSafe(data []byte) string {
	return base58.Encode(data)
}

// DecodeURLSafe reverts EncodeURLSafe.
func DecodeURLSafe(data string) ([]byte, error) {
	return base58.Decode(data)
}
