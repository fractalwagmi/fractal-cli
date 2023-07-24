package crc32c_test

import (
	b64 "encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/fractalwagmi/fractal-cli/pkg/crc32c"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCrc32c(t *testing.T) {
	tempDir := t.TempDir()
	localFilename := filepath.Join(tempDir, uuid.NewString()+".png")
	downloadFile(t, "https://storage.googleapis.com/fractal-game-releases-test/dog_that_goes_everywhere.png", localFilename)

	hash, err := crc32c.GenerateCrc32C(localFilename)
	require.NoError(t, err)

	assert.Equal(t, "C6DjLg==", b64.StdEncoding.EncodeToString(hash))
}

func downloadFile(t *testing.T, inputUrl string, destFile string) {
	t.Logf("downloading: %s to %s", inputUrl, destFile)

	out, err := os.Create(destFile)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Get(inputUrl)
	require.NoError(t, err)

	require.Equal(t, resp.StatusCode, http.StatusOK)

	if _, err := io.Copy(out, resp.Body); err != nil {
		t.Fatal(err)
	}

	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
}
