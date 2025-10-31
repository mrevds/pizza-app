.PHONY: help build run test clean docker-build docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - Show Docker logs"

build:
	go build -o bin/api-gateway cmd/api-gateway/main.go

run: build
	./bin/api-gateway

test:
	go test ./...

clean:
	rm -rf bin/
	go clean

docker-build:
	docker build -t api-gateway:latest .

docker-up:
	docker-compose -f ../docker-compose.yaml up -d api-gateway

docker-down:
	docker-compose -f ../docker-compose.yaml down

docker-logs:
	docker-compose -f ../docker-compose.yaml logs -f api-gateway

deps:
	go mod download
	go mod tidy

