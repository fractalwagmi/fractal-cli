package auth

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GenerateToken(clientId string, clientSecret string) (string, error) {
	url := "https://auth-api.fractal.is/auth/oauth/token"

	payload := strings.NewReader(fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\"}", clientId, clientSecret))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return string(body), nil
}
