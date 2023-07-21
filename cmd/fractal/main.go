package main

import (
	"flag"
	"log"
	"strings"

	"github.com/fractalwagmi/fractal-cli/pkg/auth"
	"github.com/fractalwagmi/fractal-cli/pkg/crc32c"
	"github.com/fractalwagmi/fractal-cli/pkg/sdk"
	"github.com/fractalwagmi/fractal-cli/pkg/storage"
)

type Args struct {
	archive            string
	clientId           string
	clientSecret       string
	displayName        string
	platform           string
	version            string
	exeFile            string
	macAppDirectory    string
	macInnerExecutable string
}

var platforms = map[string]bool{
	"WINDOWS":   true,
	"MAC":       true,
	"UNIVERSAL": true,
}

func main() {
	args := Args{}

	flag.StringVar(&args.archive, "zip", "", "path to .zip archive of game binary")
	flag.StringVar(&args.clientId, "clientId", "", "Fractal client id")
	flag.StringVar(&args.clientSecret, "clientSecret", "", "Fractal client secret")
	flag.StringVar(&args.displayName, "displayName", "", "Display name for the build file (optional, not shown to end users)")
	flag.StringVar(&args.platform, "platform", "", "Binary platform (windows, mac, universal)")
	flag.StringVar(&args.version, "version", "", "Binary version (must be unique to project)")
	flag.StringVar(&args.exeFile, "exeFile", "", "Path to .exe file to run on launch (windows only, required)")
	flag.StringVar(&args.macAppDirectory, "macAppDirectory", "", "Path to macOS .app directory (mac only, required)")
	flag.StringVar(&args.macInnerExecutable, "macInnerExecutable", "", "Path to macOS inner executable (mac only, required). Example: game.app/Contents/MacOS/game")

	flag.Parse()

	sanitizeArguments(&args)
	validateArguments(args)

	if token, err := auth.GenerateToken(args.clientId, args.clientSecret); err != nil {
		log.Fatalf("Error generating token: %s\n", err)
	} else if crc32c, err := crc32c.GenerateCrc32C(args.archive); err != nil {
		log.Fatalf("Error generating crc32c: %s\n", err)
	} else if res, err := sdk.CreateBuild(token, crc32c, args.displayName); err != nil {
		log.Fatalf("Error creating build: %s\n", err)
	} else if err := storage.UploadFile(res.UploadUrl, args.archive); err != nil {
		log.Fatalf("Error uploading file: %s\n", err)
	} else if err := sdk.UpdateBuild(token, sdk.UpdateBuildRequest{
		BuildNumber:        res.BuildNumber,
		Platform:           args.platform,
		Version:            args.version,
		ExeFile:            args.exeFile,
		MacAppDirectory:    args.macAppDirectory,
		MacInnerExecutable: args.macInnerExecutable,
	}); err != nil {
		log.Fatalf("Error updating build: %s\n", err)
	}
}

func sanitizeArguments(args *Args) {
	args.platform = strings.ToUpper(args.platform)

	// display name is optional, so set fallback from file name if not provided
	if args.displayName == "" {
		args.displayName = args.archive[strings.LastIndex(args.archive, "/")+1:]
	}
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

	// validate platform
	if !platforms[args.platform] {
		log.Fatalf("Invalid platform: %s\n", args.platform)
	}

	// required fields
	if args.platform == "WINDOWS" || args.platform == "UNIVERSAL" {
		if len(args.exeFile) == 0 {
			log.Fatal("Missing exeFile\n")
		}
	}
	if args.platform == "MAC" || args.platform == "UNIVERSAL" {
		if len(args.macAppDirectory) == 0 {
			log.Fatal("Missing mac app directory\n")
		}
		if len(args.macInnerExecutable) == 0 {
			log.Fatal("Missing mac inner executable\n")
		}
	}

	// excluded fields
	if args.platform == "WINDOWS" && (len(args.macAppDirectory) > 0 || len(args.macInnerExecutable) > 0) {
		log.Fatal("Mac fields (app directory or inner executable) were defined for a windows build.\n")
	}
	if args.platform == "MAC" && len(args.exeFile) > 0 {
		log.Fatal("Windows fields were defined for a mac build.\n")
	}
}
