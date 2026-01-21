.PHONY: all build clean test install help run dry-run list-reservations

# Variables
BINARY_NAME=oci-insta-scale
CAPACITY_MANAGER=capacity-manager
CONFIG_FILE=config.yaml

# Build both binaries
all: build

# Build all binaries
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) main.go
	@echo "Building $(CAPACITY_MANAGER)..."
	@go build -o $(CAPACITY_MANAGER) capacity-manager.go
	@echo "✓ Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(CAPACITY_MANAGER)
	@echo "✓ Clean complete!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies ready!"

# Run tests (if any)
test:
	@echo "Running tests..."
	@go test -v ./...

# Install binaries to GOPATH/bin
install: build
	@echo "Installing binaries to GOPATH/bin..."
	@cp $(BINARY_NAME) $(GOPATH)/bin/
	@cp $(CAPACITY_MANAGER) $(GOPATH)/bin/
	@echo "✓ Installation complete!"

# Run with default config
run: build
	@./$(BINARY_NAME)

# Run in dry-run mode
dry-run: build
	@./$(BINARY_NAME) -dry-run

# List capacity reservations
list-reservations: build
	@./$(CAPACITY_MANAGER) -list

# Create example reservation
create-reservation: build
	@./$(CAPACITY_MANAGER) -create \
		-name "example-reservation" \
		-ad "rgiR:US-ASHBURN-AD-1" \
		-shape "VM.Standard.E4.Flex" \
		-count 5 \
		-ocpus 1 \
		-memory 6

# Setup config from example
setup-config:
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "Creating $(CONFIG_FILE) from example..."; \
		cp config.example.yaml $(CONFIG_FILE); \
		echo "✓ Config file created. Please edit $(CONFIG_FILE) with your settings."; \
	else \
		echo "$(CONFIG_FILE) already exists. Not overwriting."; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete!"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "golangci-lint not installed. Run: brew install golangci-lint"

# Show help
help:
	@echo "Available commands:"
	@echo "  make build              - Build all binaries"
	@echo "  make clean              - Remove build artifacts"
	@echo "  make deps               - Download Go dependencies"
	@echo "  make test               - Run tests"
	@echo "  make install            - Install binaries to GOPATH/bin"
	@echo "  make run                - Build and run with default config"
	@echo "  make dry-run            - Run in dry-run mode (no instances created)"
	@echo "  make list-reservations  - List capacity reservations"
	@echo "  make create-reservation - Create example capacity reservation"
	@echo "  make setup-config       - Create config.yaml from example"
	@echo "  make fmt                - Format Go code"
	@echo "  make lint               - Lint Go code"
	@echo "  make help               - Show this help message"
