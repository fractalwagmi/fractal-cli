package storage_test

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	gstorage "cloud.google.com/go/storage"
	"github.com/fractalwagmi/fractal-cli/pkg/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUploadFile(t *testing.T) {
	fakeGcs, error := fakestorage.NewServerWithOptions(fakestorage.Options{
		Host:       "localhost",
		Port:       8080,
		Scheme:     "http",
		PublicHost: "127.0.0.1:8080",
	})
	require.NoError(t, error)
	defer fakeGcs.Stop()

	storageClient := fakeGcs.Client()

	// ensure that bucket is created
	if err := storageClient.Bucket("test-bucket").Create(context.Background(), "fake-project", nil); err != nil {
		t.Fatal(err)
	}

	uploadUrl, err := generateSignedUploadUrl(t, storageClient, "test-bucket", "test-object")
	require.NoError(t, err)

	z := generateZip(t, "test file")

	decodedUrl, error := url.QueryUnescape(uploadUrl)
	require.NoError(t, error)

	ru, err := storage.GetResumableUploadUrl(http.DefaultClient, decodedUrl)
	require.NoError(t, err)

	require.NoError(t, storage.UploadFile(http.DefaultClient, ru, z))
}

func generateSignedUploadUrl(
	t *testing.T,
	storageClient *gstorage.Client,
	bucket string,
	object string,
) (string, error) {
	r, err := os.Open("testing/fake-but-valid-key.pem")
	require.NoError(t, err)

	pk, err := io.ReadAll(r)
	require.NoError(t, err)

	opts := &gstorage.SignedURLOptions{
		GoogleAccessID: "test@serviceaccount.com",
		PrivateKey:     pk,
		Scheme:         gstorage.SigningSchemeV4,
		ContentType:    "application/zip",
		Method:         "POST",
		Insecure:       true,
		Expires:        time.Now().Add(5 * time.Minute),
		Style:          gstorage.BucketBoundHostname("127.0.0.1:8080/" + bucket),
		Headers: []string{
			"x-goog-resumable:start",
			"Content-Length:0",
		},
	}

	if signedUrl, err := storageClient.Bucket(bucket).SignedURL(object, opts); err != nil {
		return "", err
	} else {
		return signedUrl, nil
	}
}

func generateZip(t *testing.T, contents string) string {
	filename := t.TempDir() + "/" + uuid.NewString() + ".zip"
	file, err := os.Create(filename)
	require.NoError(t, err)

	// Create a new zip writer
	wr := zip.NewWriter(file)

	// Add a file to the zip file
	f, err := wr.Create("test.txt")
	require.NoError(t, err)

	// Write data to the file
	_, err = f.Write([]byte(contents))
	require.NoError(t, err)

	if err := wr.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	return filename
}
