package main

import (
	"strings"

	"github.com/fractalwagmi/fractal-cli/internal/upload"

	"log"
	"os"
)

const (
	uploadCommandName = "upload"
)

var commandMap = map[string]bool{
	uploadCommandName: true,
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf(
			"Please provide a command. Available commands: [%s]\n",
			strings.Join(commands(), ", "),
		)
	}

	command := os.Args[1]

	checkForHelp(command)
	checkValidCommand(command)

	os.Args = os.Args[1:]

	if command == uploadCommandName {
		upload.Run()
	}
}

func commands() []string {
	var commands []string
	for c := range commandMap {
		commands = append(commands, c)
	}
	return commands
}

func checkForHelp(command string) {
	if command == "-h" || command == "-help" || command == "--help" || command == "--h" {
		log.Fatalf(
			"To learn more about usage, please invoke one of the following commands with -h: [%s]\n",
			strings.Join(commands(), ", "),
		)
	}
}

func checkValidCommand(command string) {
	if !commandMap[command] {
		log.Fatalf(
			"Invalid command: %s. Please use one of the following commands: [%s]\n",
			command,
			strings.Join(commands(), ", "),
		)
	}
}
