package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"

	"github.com/leclerc04/go-tool/agl/util/strs"

	"github.com/leclerc04/go-tool/agl/util/errs"

	"github.com/leclerc04/go-tool/agl/util/must"
)

// Cryptor encrypts/decrypts data.
type Cryptor struct {
	cipher cipher.AEAD
}

// NewCryptor creates a cryptor.
// keyStr must be a hex encoded 128/256 bit key. Generate it at
// https://keygen.io/ and choose the SHA 128 key.
func NewCryptor(keyStr string) *Cryptor {
	key, err := hex.DecodeString(keyStr)
	must.Must(err)
	block, err := aes.NewCipher(key)
	must.Must(err)
	aesgcm, err := cipher.NewGCM(block)
	must.Must(err)
	return &Cryptor{
		cipher: aesgcm,
	}
}

// Encrypt encrypts the plaintext to a url safe string.
func (c *Cryptor) Encrypt(plaintext []byte) string {
	nonce := make([]byte, c.cipher.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	must.Must(err)
	ciphertext := c.cipher.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(ciphertext, nonce...)
	return "v1v" + strs.EncodeURLSafe(ciphertext)
}

// Decrypt decrypts the ciphertext.
func (c *Cryptor) Decrypt(ciphertextStr string) ([]byte, error) {
	ciphertextStr = strings.TrimPrefix(ciphertextStr, "v1v")
	ciphertext, err := strs.DecodeURLSafe(ciphertextStr)
	if err != nil {
		return nil, errs.InvalidArgument.Wrapf(err, "invalid data")
	}
	nonce := ciphertext[len(ciphertext)-c.cipher.NonceSize():]
	ciphertext = ciphertext[:len(ciphertext)-c.cipher.NonceSize()]
	plaintext, err := c.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errs.InvalidArgument.Wrapf(err, "invalid data")
	}
	return plaintext, nil
}
