package sdk

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseUrl = "https://api.fractal.is/sdk"

type CreateBuildRequest struct {
	DisplayName string `json:"display_name"`
	Crc32C      string `json:"crc32c"`
}

type CreateBuildResponse struct {
	UploadUrl   string `json:"uploadUrl"`
	BuildNumber uint32 `json:"buildNumber"`
}

// Creates a build and returns an upload URL to upload the binary
func CreateBuild(
	authToken string,
	crc32c []byte,
	displayName string,
) (*CreateBuildResponse, error) {
	url := baseUrl + "/launcher/build/create"

	reqBody := CreateBuildRequest{
		DisplayName: displayName,
		Crc32C:      hex.EncodeToString(crc32c),
	}
	bodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	out := &CreateBuildResponse{}
	if err := json.Unmarshal(body, out); err != nil {
		return nil, err
	}

	return out, nil
}

type UpdateBuildRequest struct {
	BuildNumber        uint32 `json:"build_number"`
	Platform           string `json:"platform"`
	Version            string `json:"version"`
	ExeFile            string `json:"exe_file,omitempty"`
	MacAppDirectory    string `json:"mac_app_directory,omitempty"`
	MacInnerExecutable string `json:"mac_inner_executable,omitempty"`
}

// Creates a build and returns an upload URL to upload the binary
func UpdateBuild(
	authToken string,
	update UpdateBuildRequest,
) error {
	// TODO(john): remove this and wait for some signal from API when zip
	// processing is complete. Otherwise, we will get an error saying that windows
	// or mac files were not found in the .zip archive.
	time.Sleep(5 * time.Second)

	url := baseUrl + "/launcher/build/" + fmt.Sprint(update.BuildNumber) + "/update"

	body, err := json.Marshal(update)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	if res, err := http.DefaultClient.Do(req); err != nil {
		return err
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
