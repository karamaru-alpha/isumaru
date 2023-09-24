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
RUN go build -o web cmd/web/main.go

# ==== build docker image ====
FROM alpine:${ALPINE_VERSION}

WORKDIR /usr/src/karamaru-alpha/isumaru

RUN wget https://github.com/tkuchiki/alp/releases/download/v1.0.16/alp_linux_amd64.zip
RUN unzip alp_linux_amd64.zip
RUN mv alp /usr/local/bin/

RUN wget https://github.com/tkuchiki/slp/releases/download/v0.1.2/slp_linux_amd64.zip
RUN unzip slp_linux_amd64.zip
RUN mv slp /usr/local/bin/

COPY --from=go-builder /go/src/github.com/karamaru-alpha/isumaru/isumaru isumaru

RUN chmod +x /usr/src/karamaru-alpha/web/web
ENTRYPOINT ["/usr/src/karamaru-alpha/isumaru/isumaru"]
