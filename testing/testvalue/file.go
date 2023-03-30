package testvalue

import (
	"os"
)

// MustReadFile reads the file and returns the content byte slice.
// nolint:gosec
func MustReadFile(filePath string) []byte {
	b, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return b
}
