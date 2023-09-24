include .env

.PHONY: run
run:
	PORT=${ISUMARU_PORT} AGENT_URL=${AGENT_URL} go run cmd/isumaru/main.go

.PHONY: run-agent
run-agent:
	PORT=${AGENT_PORT} go run cmd/agent/main.go

.PHONY: run-web
run-web:
	(cd web; pnpm run dev)

.PHONY: access
access:
	@printf "# Time: 2022-01-07T07:54:07.943312Z\n# User@Host: root[root] @ localhost []  Id:    10\n# Query_time: 0.410568  Lock_time: 0.000513 Rows_sent: 1  Rows_examined: 2844047\nSET timestamp=1641542047;\nselect count(*) from salaries where salary >= 10000;\n" >> testdata/slow-query.log
	@echo "time:2015-09-06T06:00:43+09:00	method:GET	uri:/diary/entry/5678	status:200	size:30	apptime:0.432" >> testdata/access.log
