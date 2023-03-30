package testvalue

import (
	"math/rand"
	"time"
)

// NotNilRnd returns a not nil random generator: the source or creates a new.
func NotNilRnd(rnd *rand.Rand) *rand.Rand {
	if rnd != nil {
		return rnd
	}

	// nolint:gosec
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
