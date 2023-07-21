package storage

func UploadFile(url string, file string) error {
	_, err := getResumableUploadUrl(url)
	if err != nil {
		return err
	}

	// TODO(john): implement actual upload.

	return nil
}

func getResumableUploadUrl(url string) (string, error) {
	// TODO(john): implement (exchange api-provided url for gcp 'location')
	return "", nil
}
