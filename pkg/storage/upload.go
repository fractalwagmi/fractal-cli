package storage

import (
	"net/http"
	"os"
)

func UploadFile(url string, file string) error {
	uploadUrl, err := getResumableUploadUrl(url)
	if err != nil {
		return err
	}

	// TODO(john): upload in chunks and with parallelism to increase upload speed.
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

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func getResumableUploadUrl(url string) (string, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Length", "0")
	req.Header.Add("Content-Type", "application/zip")
	req.Header.Add("x-goog-resumable", "start")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	return res.Header.Get("Location"), nil
}
