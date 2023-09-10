ARG GO_VERSION=1.21.1
ARG ALPINE_VERSION=3.18

# ===== build go binary =====
FROM golang:${GO_VERSION}-alpine as go-builder

WORKDIR /go/src/github.com/karamaru-alpha/isumaru

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg
RUN go build -o isumaru-agent cmd/agent/main.go

# ==== build docker image ====
FROM alpine:${ALPINE_VERSION}

WORKDIR /usr/src/karamaru-alpha/isumaru

COPY --from=go-builder /go/src/github.com/karamaru-alpha/isumaru/isumaru-agent isumaru-agent

RUN chmod +x /usr/src/karamaru-alpha/isumaru/isumaru-agent
ENTRYPOINT ["/usr/src/karamaru-alpha/isumaru/isumaru-agent"]
