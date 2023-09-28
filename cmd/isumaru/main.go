package main

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/isumaru/pkg/isumaru"
)

func main() {
	os.Exit(cmd())
}

func cmd() (code int) {
	config := &isumaru.Config{
		Port:                  getEnv("PORT", "8000"),
		SlowQueryLogDirFormat: getEnv("SLOW_QUERY_LOG_DIR_FORMAT", "log/%s/slowquerylog"),
		AccessLogDirFormat:    getEnv("ACCESS_LOG_DIR_FORMAT", "log/%s/accesslog"),
		PProfDirFormat:        getEnv("PPROF_LOG_DIR_FORMAT", "log/%s/pprof"),
		SlpConfigPath:         getEnv("SLP_CONFIG_PATH", "config/slp.yaml"),
		AlpConfigPath:         getEnv("ALP_CONFIG_PATH", "config/alp.yaml"),
	}

	isumaru.Serve(config)
	return 0
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	slog.Error("Env value %s is not found", key)
	return fallback
}
