package testing

import (
	"os"
	"testing"
)

// Alternative to testing.T.TempDir() that does not fail when RemoveAll fails (cuz RemoveAll doesnt work on windows)
// https://github.com/golang/go/issues/51442
func TempDir(t *testing.T) string {
	dir := os.TempDir()
	t.Cleanup(func() {
		// this is best effort
		// https://github.com/golang/go/issues/51442
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("unable to delete: %s err: %v", dir, err)
		}
	})
	return dir
}
