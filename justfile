bin := "dist"

default: clean build-cli build-server

build-cli:
    go build -ldflags="-s -w" -o {{bin}}/cli ./cmd/cli

install-cli: build-cli
    ./scripts/install.sh -local bin/cli

build-server:
    go build -ldflags="-s -w" -o {{bin}}/server ./cmd/server

server:
    go run ./cmd/server

cli *args:
    go run ./cmd/cli {{args}}

dev:
    docker compose up --watch --build

sudo-dev:
    sudo docker compose up --watch --build

clean:
    rm -rf {{bin}}
