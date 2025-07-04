package randutil

import (
	"math/rand"
	"sync"
)

var seedOnce sync.Once

// Seed the weak global rand in case we forgot.
func init() {
	SeedGlobalMathRand()
}

// SeedGlobalMathRand seeds the math/rand. It guarantees not
// seeding more than once per process. (Seeding arbitrarily is dangerous).
// Note that even though it is seeded with crypto source, the
// rand is still not safe for crypto purpose.
func SeedGlobalMathRand() {
	seedOnce.Do(func() {
		rand.Seed(Rand.Int63())
	})
}
