# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Wipe everything and start fresh
docker-reset:
	docker-compose down --volumes

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Migration commands
migrate-up:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go -command=up

migrate-down:
	@echo "Rolling back last migration..."
	@go run cmd/migrate/main.go -command=down

migrate-status:
	@echo "Checking migration status..."
	@go run cmd/migrate/main.go -command=status

migrate-reset:
	@echo "Resetting database..."
	@go run cmd/migrate/main.go -command=reset

# Create a new migration file
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	filename="internal/migration/goose/$${timestamp}_$(name).sql"; \
	echo "-- +goose Up" > $$filename; \
	echo "-- Add your migration SQL here" >> $$filename; \
	echo "" >> $$filename; \
	echo "-- +goose Down" >> $$filename; \
	echo "-- Add your rollback SQL here" >> $$filename; \
	echo "Created migration file: $$filename"

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest migrate-up migrate-down migrate-status migrate-reset migrate-create
