BIN=bin

all: clean build-cli build-server

build-cli:
	go build -ldflags="-s -w" -o $(BIN)/cli ./cmd/cli

install-cli:
	make build-cli
	./scripts/install.sh -local bin/cli

build-server:
	go build -ldflags="-s -w" -o $(BIN)/server ./cmd/server

run-server:
	go run ./cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf $(BIN)
