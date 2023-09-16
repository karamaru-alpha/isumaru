package main

import (
	"os"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru"
)

func main() {
	os.Exit(cmd())
}

func cmd() (code int) {
	config := &isumaru.Config{
		Port:     getEnv("PORT", "8000"),
		AgentURL: getEnv("AgentURL", "http://localhost:19000"),
	}

	isumaru.Serve(config)
	return 0
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
