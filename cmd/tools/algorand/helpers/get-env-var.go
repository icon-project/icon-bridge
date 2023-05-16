package helpers

import (
	"log"
	"os"
)

func GetEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s environment variable is not set", name)
	}
	return value
}