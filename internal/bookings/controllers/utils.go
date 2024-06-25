package controllers

import "fmt"

func paramRequired(name string) string {
	return fmt.Sprintf("param '%v' is required", name)
}

func deserializationError(err error) string {
	return fmt.Sprintf("unable to deserialize value: %s", err)
}

func serializationError(err error) string {
	return fmt.Sprintf("unable to serialize value: %s", err)
}

func tariffError(err error) string {
	return fmt.Sprintf("unable to get basic tariff: %s", err)
}
