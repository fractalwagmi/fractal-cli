package storage

import (
	"net/http"
	"os"
)

func UploadFile(httpClient *http.Client, uploadUrl string, file string) error {
	// TODO(john): upload in chunks to protect against transient network errors.
	data, err := os.Open(file)
	if err != nil {
		return err
	}
	defer data.Close()
	req, err := http.NewRequest("PUT", uploadUrl, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/zip")

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

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

	defer res.Body.Close()

	return res.Header.Get("Location"), nil
}
