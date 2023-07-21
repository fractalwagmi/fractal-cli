package sdk

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
)

const baseUrl = "https://api.fractal.is/sdk"

type CreateBuildRequest struct {
	DisplayName string `json:"display_name"`
	Crc32C      string `json:"crc32c"`
}

type CreateBuildResponse struct {
	UploadUrl   string `json:"upload_url"`
	BuildNumber uint32 `json:"build_number"`
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
