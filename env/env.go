package env

import (
	"log"
	"os"
)

func MandatoryString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("you must define %s environment variable", key)
	}
	return val
}
