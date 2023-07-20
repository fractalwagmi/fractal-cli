package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://api.fractal.is/sdk"

type CreateBuildResponse struct {
	UploadUrl   string `json:"upload_url"`
	BuildNumber uint32 `json:"build_umber"`
}

// Creates a build and returns an upload URL to upload the binary
func CreateBuild(authToken string) (*CreateBuildResponse, error) {
	url := baseUrl + "/launcher/build/create"

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	out := &CreateBuildResponse{}
	if err := json.Unmarshal(body, out); err != nil {
		return nil, err
	}

	return out, nil
}
