package storage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultChunkSize = 1 << 22 // 4MB in bytes
)

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

	println("Uploading file...")

	// Upload file in chunks with retry policy for resilience against transient connectivity errors.
	for offset := int64(0); offset < fileSize; offset += defaultChunkSize {
		chunk := make([]byte, defaultChunkSize)
		n, err := file.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}

		if err := uploadChunkWithRetryPolicy(
			http.DefaultClient,
			uploadUrl,
			// Only use the part of the chunk that contains data (n bytes).
			chunk[:n],
			offset,
			fileSize,
		); err != nil {
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

var backoffSchedule = []time.Duration{
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
}

func uploadChunkWithRetryPolicy(
	client *http.Client,
	uploadUrl string,
	chunk []byte,
	offset int64,
	fileSize int64,
) error {
	var err error

	for _, backoff := range backoffSchedule {
		err = uploadChunk(client, uploadUrl, chunk, offset, fileSize)
		if err == nil {
			break
		}

		time.Sleep(backoff)
	}

	return err
}

func uploadChunk(
	client *http.Client,
	uploadUrl string,
	chunk []byte,
	offset int64,
	fileSize int64,
) error {
	req, err := http.NewRequest("PUT", uploadUrl, bytes.NewReader(chunk))
	if err != nil {
		return err
	}

	chunkSize := int64(len(chunk))
	byteRange := fmt.Sprintf("bytes %d-%d/%d", offset, offset+chunkSize-1, fileSize)
	lastChunk := offset+defaultChunkSize >= fileSize

	req.Header.Set("Content-Length", fmt.Sprint(chunkSize))
	req.Header.Set("Content-Range", byteRange)

	res, err := http.DefaultClient.Do(req)
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
