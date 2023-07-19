package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/fractalwagmi/fractal-cli/pkg/auth"
)

func main() {
	archive := flag.String("zip", "", "path to .zip archive of game binary")
	clientId := flag.String("clientId", "", "Fractal client id")
	clientSecret := flag.String("clientSecret", "", "Fractal client secret")

	flag.Parse()

	validateArguments(archive, clientId, clientSecret)

	token, err := auth.GenerateToken(*clientId, *clientSecret)
	if err != nil {
		log.Fatalf("Error generating token: %s\n", err)
	}

	fmt.Printf("Token: %s\n", token)
}

func validateArguments(
	archive *string,
	clientId *string,
	clientSecret *string,
) {
	if !strings.HasSuffix(*archive, ".zip") {
		log.Fatalf("Invalid archive: %s\n", *archive)
	}
	if len(*clientId) == 0 {
		log.Fatalf("Invalid clientId: %s\n", *clientId)
	}
	if len(*clientSecret) == 0 {
		log.Fatalf("Invalid clientSecret: %s\n", *clientSecret)
	}
}
