.PHONY: all build test lint fmt tidy clean

MODULES := ./docker-wrapper

all: build

## build: build all modules
build:
	go build $(MODULES)/...

## test: run tests for all modules
test:
	go test $(MODULES)/...

## lint: run go vet on all modules
lint:
	go vet $(MODULES)/...

## fmt: format all Go source files
fmt:
	gofmt -w $(MODULES)

## tidy: tidy dependencies for all modules
tidy:
	go work sync
	@for mod in $(MODULES); do \
		echo ">> tidy $$mod"; \
		cd $$mod && go mod tidy && cd -; \
	done

## clean: remove build artifacts
clean:
	go clean $(MODULES)/...

## help: print this help message
help:
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
