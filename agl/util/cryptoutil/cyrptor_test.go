package cryptoutil_test

import (
	"testing"

	"github.com/leclecr04/go-tool/agl/util/cryptoutil"
	"github.com/stretchr/testify/assert"
)

func TestCryptor(t *testing.T) {
	c := cryptoutil.NewCryptor("93C13D5AAAF33F1593C13D5AAAF33F15")
	ctext := c.Encrypt([]byte("hello world"))
	ptext, err := c.Decrypt(ctext)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", string(ptext))
}
