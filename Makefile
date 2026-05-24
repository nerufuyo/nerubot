# Makefile for NeruBot (Rust Edition)

# Variables
BINARY_NAME=nerubot
CARGO=cargo

.PHONY: all build clean test fmt lint run help docker-build docker-run docker-stop check

## all: Default target - build the application
all: build

## build: Build release binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(CARGO) build --release
	@echo "Build complete: target/release/$(BINARY_NAME)"

## run: Build and run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(CARGO) run

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(CARGO) clean
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	$(CARGO) test

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(CARGO) fmt
	@echo "Format complete"

## lint: Run clippy linter
lint:
	@echo "Running clippy..."
	$(CARGO) clippy -- -D warnings

## check: Run typecheck (cargo check)
check:
	@echo "Checking..."
	$(CARGO) check
	@echo "Check complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .
	@echo "Docker image built: $(BINARY_NAME):latest"

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker-compose up -d
	@echo "Container started"

## docker-stop: Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker-compose down
	@echo "Container stopped"

## all-checks: Run all checks (fmt, check, lint, test)
all-checks: fmt check lint test
	@echo "All checks passed!"

## help: Show this help message
help:
	@echo "NeruBot Makefile Commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
