package storage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fractalwagmi/fractal-cli/pkg/functions"
	progressbar "github.com/schollz/progressbar/v3"
)

const (
	defaultChunkSize = 1 << 22 // 4MB in bytes
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
}

func UploadFile(httpClient *http.Client, uploadUrl string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file's size.
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	bar := progressbar.DefaultBytes(fileSize, "Uploading game binary")

	// Upload file in chunks with retry policy for resilience against transient connectivity errors.
	for offset := int64(0); offset < fileSize; offset += defaultChunkSize {
		chunk := make([]byte, defaultChunkSize)
		n, err := file.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}

		if err := functions.RunWithRetryPolicy(
			backoffSchedule,
			func() error {
				// chunks may be uploaded multiple times with retries, so always reset to
				// 'offset' at the beginning on the individual chunk upload.
				if err := bar.Set64(offset); err != nil {
					return err
				}

				// Only use the part of the chunk that contains data (n bytes).
				return uploadChunk(httpClient, uploadUrl, chunk[:n], offset, fileSize, bar)
			}); err != nil {
			return err
		}
	}

	println("File upload completed successfully!")
	return nil
}

func GetResumableUploadUrl(httpClient *http.Client, url string) (string, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Length", "0")
	req.Header.Add("Content-Type", "application/zip")
	req.Header.Add("x-goog-resumable", "start")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	return res.Header.Get("Location"), nil
}

type wrapper struct {
	io.Reader
	bar *progressbar.ProgressBar
	n   int
}

func (w *wrapper) Read(p []byte) (int, error) {
	n, err := w.Reader.Read(p)
	w.n += n
	w.bar.Add(n)
	return n, err
}

func uploadChunk(
	client *http.Client,
	uploadUrl string,
	chunk []byte,
	offset int64,
	fileSize int64,
	bar *progressbar.ProgressBar,
) error {
	w := &wrapper{
		Reader: bytes.NewReader(chunk),
		bar:    bar,
	}

	req, err := http.NewRequest("PUT", uploadUrl, w)
	if err != nil {
		return err
	}

	chunkSize := int64(len(chunk))
	byteRange := fmt.Sprintf("bytes %d-%d/%d", offset, offset+chunkSize-1, fileSize)
	lastChunk := offset+defaultChunkSize >= fileSize

	req.Header.Set("Content-Length", fmt.Sprint(chunkSize))
	req.Header.Set("Content-Range", byteRange)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// 308 status code for successful request on all chunks except last chunk
	if !lastChunk && res.StatusCode != 308 {
		return fmt.Errorf("chunk upload failed: %v", err)
	}

	if lastChunk && res.StatusCode != http.StatusOK {
		return fmt.Errorf("chunk upload failed: %v", err)
	}

	return nil
}
