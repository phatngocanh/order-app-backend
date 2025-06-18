package env

import (
	"errors"
	"os"
)

func GetEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}

	return "", errors.New("environment variable not found")
}
