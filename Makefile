MAKEFLAGS += --silent

build:
	go build -o ./tmp/subscriptions ./cmd/subscriptions/main.go

build-debug:
	go build -o ./tmp/subscriptions -gcflags="all=-N -l" ./cmd/exchange/main.go

run:
	go run ./cmd/subscriptions/main.go
