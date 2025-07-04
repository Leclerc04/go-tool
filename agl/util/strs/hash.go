package strs

import (
	"crypto/md5"
	"fmt"
)

// MD5Sum computes the md5 for a string.
func MD5Sum(a string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(a)))
}
