package helper

import (
	"fmt"
	"os"
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
