# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o subscriptions ./cmd/subscriptions
RUN chmod +x ./subscriptions

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache make curl
# Install dockerize
RUN curl -L https://github.com/jwilder/dockerize/releases/download/v0.9.3/dockerize-alpine-linux-amd64-v0.9.3.tar.gz | tar xzf - -C /usr/local/bin && chmod +x /usr/local/bin/dockerize
# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin && chmod +x /usr/local/bin/migrate
COPY --from=builder /app/subscriptions ./
COPY --from=builder /app/Makefile ./
