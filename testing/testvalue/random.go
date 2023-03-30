package testvalue

import (
	"math/rand"
	"time"
)

// NotNilRnd returns a not nil random generator: the source or creates a new.
func NotNilRnd(r *rand.Rand) *rand.Rand {
	if r == nil {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return r
}
