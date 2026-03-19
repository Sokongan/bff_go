package config

import (
	"fmt"
	"os"
	"strings"
)

const defaultHTTPAddress = ":8080"

func loadServerAddress() string {
	if addr := strings.TrimSpace(os.Getenv("HTTP_ADDR")); addr != "" {
		return addr
	}

	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.HasPrefix(port, ":") {
			return port
		}
		return fmt.Sprintf(":%s", port)
	}

	return defaultHTTPAddress
}
