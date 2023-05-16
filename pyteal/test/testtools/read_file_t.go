package testtools

import (
	"os"
	"testing"
)

func ReadFileT(t *testing.T, path string) (content []byte) {
	t.Helper()

	var err error
	content, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %s\n", err)
	}
	return
}
