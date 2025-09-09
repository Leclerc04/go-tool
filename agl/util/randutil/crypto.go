package randutil

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"

	"github.com/leclerc04/go-tool/agl/util/must"
)

// Crypto is a cryptographically strong Rand. Use it for tokens, password etc.
var Crypto = rand.New(&cryptoSrc{})

// Rand is the default rand that should be used.
var Rand = Crypto

type cryptoSrc struct{}

func (c *cryptoSrc) Seed(seed int64) {
	// No-op.
}

func (c *cryptoSrc) Uint64() (value uint64) {
	must.Must(binary.Read(crand.Reader, binary.BigEndian, &value))
	return value
}

func (c *cryptoSrc) Int63() int64 {
	return int64(c.Uint64() & ^uint64(1<<63))
}
