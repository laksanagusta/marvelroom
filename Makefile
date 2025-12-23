# Variables
APP_NAME := reika
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "latest")
REGISTRY := laksanadika
IMAGE_NAME := $(REGISTRY)/$(APP_NAME)

# Docker commands
.PHONY: build-image
build-image:
	docker build --no-cache --platform linux/amd64 -t $(IMAGE_NAME):$(VERSION) .
	docker tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest

.PHONY: push-image
push-image: build-image
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest
	@echo "🎉 Release complete!"

.PHONY: run-local
run-local:
	docker-compose up -d

.PHONY: stop-local
stop-local:
	docker-compose down

.PHONY: run-dev
run-dev:
	docker-compose --profile dev up

# Development commands
.PHONY: dev
dev:
	go run main.go

.PHONY: build
build:
	go build -o main .

.PHONY: test
test:
	go test -v ./...

.PHONY: tidy
tidy:
	go mod tidy

# Clean up
.PHONY: clean
clean:
	rm -f main
	docker system prune -f

# Full deployment cycle
.PHONY: deploy
deploy: test build-image push-image

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  dev              - Run application locally"
	@echo "  build            - Build Go binary"
	@echo "  test             - Run tests"
	@echo "  build-image      - Build Docker image"
	@echo "  push-image       - Push Docker image to registry"
	@echo "  run-local        - Run with Docker Compose"
	@echo "  stop-local       - Stop Docker Compose"
	@echo "  run-dev          - Run development mode with hot reload"
	@echo "  deploy           - Full deployment cycle (test → build → push)"
	@echo "  clean            - Clean up build artifacts and Docker"
	@echo "  help             - Show this help message"