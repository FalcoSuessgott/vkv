package fs

import (
	"os"
)

// ReadFile reads from a file.
func ReadFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// CreateDirectory creates a given directory recursively.
func CreateDirectory(name string) error {
	return os.MkdirAll(name, os.ModePerm)
}
