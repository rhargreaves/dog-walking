package env

import (
	"os"
)

func MustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Required environment variable not set: " + key)
	}
	return val
}
