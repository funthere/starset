package helpers

import (
	"fmt"
	"os"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	return fmt.Sprintf(":%s", port)
}
