package env

import (
	"log"
	"os"
)

func GetRequiredString(name string, alts ...string) string {
	all := append([]string{name}, alts...)
	for _, x := range all {
		v := os.Getenv(x)
		if v != "" {
			return v
		}
	}
	log.Fatalf("missing required env: %v", all)
	return ""
}
