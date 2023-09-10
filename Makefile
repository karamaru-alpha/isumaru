include .env
export

.PHONY: run-agent
run-agent:
	go run cmd/agent/main.go
