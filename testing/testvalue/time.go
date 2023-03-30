package testvalue

import "time"

// MustTime parses a time string of format RFC3339 and returns a time.
// Panics if fails.
func MustTime(in string) time.Time {
	t, err := time.Parse(time.RFC3339, in)
	if err != nil {
		panic(err)
	}

	return t
}

// MustTimePtr parses a time string of format RFC3339 and returns a pointer to a time.
// Panics if fails.
func MustTimePtr(in string) *time.Time {
	return Ptr(MustTime(in))
}

// MustTimeFn creates a time function with mocked (time.RFC3339 formatted string)time value to be returned.
// Panics if fails.
func MustTimeFn(in string) func() time.Time {
	return func() time.Time {
		return MustTime(in)
	}
}
