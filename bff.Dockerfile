# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /out/sso-bff \
	./cmd/main.go

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/sso-bff /usr/local/bin/sso-bff

ENV ENV=prod
ENV HTTP_ADDR=:8080

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/sso-bff"]
