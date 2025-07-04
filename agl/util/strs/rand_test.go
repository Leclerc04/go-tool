package strs_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/leclecr04/go-tool/agl/util/strs"
)

func TestRand(t *testing.T) {
	for i := 0; i < 20; i++ {
		rn, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			t.Error(err)
		}
		n := int(rn.Int64())
		c := strs.GenerateRandomStringURLSafe(n)
		if len(c) != n {
			t.Error(c, n)
		}
		t.Log(n)
		t.Log(c)
	}
}
