package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func GenerateToken(client *http.Client, clientId string, clientSecret string) (string, error) {
	url := "https://auth-api.fractal.is/auth/oauth/token"

	payload := strings.NewReader(fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\"}", clientId, clientSecret))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("User-Agent", "fractal-cli")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	out := &token{}
	err = json.Unmarshal(body, out)
	if err != nil {
		return "", err
	}

	return out.AccessToken, nil
}
