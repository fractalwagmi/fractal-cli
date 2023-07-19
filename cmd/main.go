package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/fractalwagmi/fractal-cli/pkg/auth"
	"github.com/fractalwagmi/fractal-cli/pkg/env"
)

func main() {
	archive := flag.String("zip", "", "path to .zip archive of game binary")

	flag.Parse()

	if !strings.HasSuffix(*archive, ".zip") {
		log.Fatalf("Invalid archive: %s\n", *archive)
	}

	clientId := env.GetRequiredString("FRACTAL_CLIENT_ID")
	clientSecret := env.GetRequiredString("FRACTAL_CLIENT_SECRET")

	token, err := auth.GenerateToken(clientId, clientSecret)
	if err != nil {
		log.Fatalf("Error generating token: %s\n", err)
	}

	fmt.Printf("Token: %s\n", token)
}
