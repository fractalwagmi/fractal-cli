package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fractalwagmi/fractal-cli/pkg/auth"
	"github.com/fractalwagmi/fractal-cli/pkg/crc32c"
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

	fr, err := os.Open(args.archive)
	if err != nil {
		log.Fatalf("Error opening archive: %s\n", err)
	}

	crc32c, err := crc32c.GenerateCrc32C(fr)
	if err != nil {
		log.Fatalf("Error generating crc32c: %s\n", err)
	}
	fr.Close()
	fmt.Printf("crc32c: %s\n", crc32c)
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
