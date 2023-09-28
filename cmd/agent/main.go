package main

import (
	"context"
	"os"

	"github.com/karamaru-alpha/isumaru/pkg/agent"
)

func main() {
	os.Exit(cmd())
}

func cmd() (code int) {
	config := &agent.Config{
		Port: getEnv("PORT", "19000"),
	}
	agent.Serve(context.Background(), config)
	return 0
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
