package randutil

import (
	"github.com/mr-tron/base58/base58"
	uuid "github.com/satori/go.uuid"
)

// NewUnique22 returns a unique id of lentgh 22.
func NewUnique22() string {
	u := uuid.NewV4()
	return base58.Encode(u.Bytes())
}
