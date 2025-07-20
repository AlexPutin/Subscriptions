MAKEFLAGS += --silent

ifneq ("$(wildcard ./.env)","")
    include .env
endif


DB_URL = "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOSTNAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

build:
	go build -o ./tmp/subscriptions ./cmd/subscriptions/main.go

build-debug:
	go build -o ./tmp/subscriptions -gcflags="all=-N -l" ./cmd/exchange/main.go

run:
	go run ./cmd/subscriptions/main.go

create-migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

migrate-up:
	migrate -path ./migrations -database $(DB_URL) up

migrate-down:
	migrate -path ./migrations -database $(DB_URL) down