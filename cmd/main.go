package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/fractalwagmi/fractal-cli/pkg/auth"
	"github.com/fractalwagmi/fractal-cli/pkg/crc32c"
	"github.com/fractalwagmi/fractal-cli/pkg/sdk"
)

type Args struct {
	archive      string
	clientId     string
	clientSecret string
}

func main() {
	args := Args{}

	flag.StringVar(&args.archive, "zip", "", "path to .zip archive of game binary")
	flag.StringVar(&args.clientId, "clientId", "", "Fractal client id")
	flag.StringVar(&args.clientSecret, "clientSecret", "", "Fractal client secret")

	flag.Parse()

	validateArguments(args)

	token, err := auth.GenerateToken(args.clientId, args.clientSecret)
	if err != nil {
		log.Fatalf("Error generating token: %s\n", err)
	}
	fmt.Printf("Auth token: %s\n", token)

	crc32c, err := crc32c.GenerateCrc32C(args.archive)
	if err != nil {
		log.Fatalf("Error generating crc32c: %s\n", err)
	}
	fmt.Printf("crc32c: %s\n", b64.StdEncoding.EncodeToString(crc32c))

	res, err := sdk.CreateBuild(token)
	if err != nil {
		log.Fatalf("Error creating build: %s\n", err)
	}

	fmt.Printf("Create build response: %v\n", res)
}

func validateArguments(args Args) {
	if !strings.HasSuffix(args.archive, ".zip") {
		log.Fatalf("Invalid archive: %s\n", args.archive)
	}
	if len(args.clientId) == 0 {
		log.Fatalf("Invalid clientId: %s\n", args.clientId)
	}
	if len(args.clientSecret) == 0 {
		log.Fatalf("Invalid clientSecret: %s\n", args.clientSecret)
	}
}
