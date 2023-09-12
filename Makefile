include .env

.PHONY: run
run:
	PORT=${ISUMARU_PORT} go run cmd/isumaru/main.go

.PHONY: run-agent
run-agent:
	PORT=${AGENT_PORT} go run cmd/agent/main.go

.PHONY: run-web
run-web:
	(cd web; pnpm run dev)
