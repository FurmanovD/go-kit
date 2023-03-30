package testvalue

// Ptr returns a pointer to an object given.
func Ptr[T any](in T) *T {
	return &in
}
