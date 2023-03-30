package testvalue

import "github.com/google/uuid"

// MustUUID parses uuid's string representation and returns a uuid.UUID.
// Panics if fails.
func MustUUID(in string) uuid.UUID {
	return uuid.MustParse(in)
}

// MustUUIDFn create uuid function with mocked uuid.
// Panics if fails.
func MustUUIDFn(in string) func() uuid.UUID {
	return func() uuid.UUID {
		return MustUUID(in)
	}
}
