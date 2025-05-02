# Makefile for domain_category_db

APP_NAME := domaindb
PKG := ./...

# Default target
.PHONY: all
all: build

## Build the project
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) main.go

## Run the application
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)

## Format all Go files
.PHONY: fmt
fmt:
	@echo "Formatting Go source files..."
	go fmt $(PKG)

## Run tests (unit + benchmarks if defined)
.PHONY: test
test:
	@echo "Running tests..."
	go test -v $(PKG)

## Run benchmarks only
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	go test -bench=. $(PKG)

## Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)

## Lint using go vet
.PHONY: lint
lint:
	@echo "Running go vet..."
	go vet $(PKG)

