BIN=dist

all: clean build-cli build-server

build-cli:
	go build -ldflags="-s -w" -o $(BIN)/cli ./cmd/cli

install-cli:
	make build-cli
	./scripts/install.sh -local bin/cli

build-server:
	go build -ldflags="-s -w" -o $(BIN)/server ./cmd/server

run:
	go run ./cmd/server

# Clean build artifacts
clean:
	rm -rf $(BIN)
