.PHONY: run build dev test clean migrate migrate-rollback migrate-rollback-all seed seed-fresh lint tidy

APP_NAME = glamping-api
BUILD_DIR = ./bin

run: build
	$(BUILD_DIR)/$(APP_NAME)

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go

dev:
	go run cmd/api/main.go

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)

migrate:
	go run cmd/migrate/main.go -action=migrate

migrate-rollback:
	go run cmd/migrate/main.go -action=rollback

migrate-rollback-all:
	go run cmd/migrate/main.go -action=rollback-all

seed:
	go run cmd/migrate/main.go -action=seed

seed-fresh:
	go run cmd/migrate/main.go -action=rollback-all
	go run cmd/migrate/main.go -action=migrate
	go run cmd/migrate/main.go -action=seed

lint:
	golangci-lint run ./...

tidy:
	go mod tidy
