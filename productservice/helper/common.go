package helper

import (
	"fmt"
	"os"
	"strings"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func ServerAddress() string {
	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":5002"
	}

	return fmt.Sprintf("%s", port)
}

func IntSliceToString(numbers []uint32) string {
	var numbersString []string
	for _, id := range numbers {
		numbersString = append(numbersString, fmt.Sprint(id))
	}
	return strings.Join(numbersString, ",")
}
