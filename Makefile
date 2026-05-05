BIN=bin

# Build all binaries
all: clean build-cli build-server

# Build the CLI
build-cli:
	@echo "Building CLI..."
	go build -ldflags="-s -w" -o $(BIN)/cli ./cmd/cli

# Build the Server
build-server:
	@echo "Building Server..."
	go build -ldflags="-s -w" -o $(BIN)/server ./cmd/server

# Run the server locally
run-server:
	go run ./cmd/server/main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BIN)
