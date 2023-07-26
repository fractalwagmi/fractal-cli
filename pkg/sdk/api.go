package sdk

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fractalwagmi/fractal-cli/pkg/functions"
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

var backoffSchedule = []time.Duration{
	500 * time.Millisecond,
	1 * time.Second,
	2 * time.Second,
	3 * time.Second,
	4 * time.Second,
	5 * time.Second,
}

func CreateBuild(
	client *http.Client,
	authToken string,
	crc32c []byte,
	displayName string,
) (*CreateBuildResponse, error) {
	println("Creating a new build...")

	var out *CreateBuildResponse

	if err := functions.RunWithRetryPolicy(backoffSchedule, func() error {
		var err error
		out, err = doCreateBuild(client, authToken, crc32c, displayName)
		return err
	}); err != nil {
		return nil, err
	} else {
		fmt.Printf("Build %d created successfully!\n", out.BuildNumber)
		return out, nil
	}
}

// Creates a build and returns an upload URL to upload the binary
func doCreateBuild(
	client *http.Client,
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
	req.Header.Add("User-Agent", "fractal-cli")

	res, err := client.Do(req)
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

func UpdateBuild(
	client *http.Client,
	authToken string,
	update UpdateBuildRequest,
) error {
	println("Updating build configuration with provided arguments...")
	if err := functions.RunWithRetryPolicy(backoffSchedule, func() error {
		return doUpdateBuild(client, authToken, update)
	}); err != nil {
		return err
	} else {
		println("Build updated successfully!")
		return nil
	}
}

// Creates a build and returns an upload URL to upload the binary
func doUpdateBuild(
	client *http.Client,
	authToken string,
	update UpdateBuildRequest,
) error {
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
	req.Header.Add("User-Agent", "fractal-cli")

	if res, err := client.Do(req); err != nil {
		return err
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
