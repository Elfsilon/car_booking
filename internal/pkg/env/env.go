package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func mustBeSpecified(name string) string {
	return fmt.Sprintf("Environment variable '%v' must be specified", name)
}

func MustLoadString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal(mustBeSpecified(key))
	}
	return value
}

func MustLoadInt(key string) int {
	value := MustLoadString(key)
	num, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return num
}
