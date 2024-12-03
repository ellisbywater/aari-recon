package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetStringNoFallback(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Println("failed to retrieve env variable")
		return "", errors.New("failed to retrieve env variable")
	}
	return val, nil
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}
