# Project variables
APP_NAME    := hello_world
BIN_DIR     := bin
MAIN_FILE   := ./cmd/main.go
PKG         := ./...
GIT_TAG     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS     := -ldflags="-s -w -X main.Version=$(GIT_TAG) -X 'main.BuildTime=$(BUILD_TIME)'"

# Go parameters
GO     ?= go
LINTER := golangci-lint

.PHONY: all run build clean lint test tidy deps docker help

all: build

## Run the application (requires .env)
run:
	@echo ">> Running $(APP_NAME)..."
	@$(GO) run $(MAIN_FILE)

## Build the binary with version injection
build:
	@echo ">> Building binary ($(GIT_TAG))..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 GOOS=linux $(GO) build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(BIN_DIR)/$(APP_NAME)"

## Clean build artifacts
clean:
	@echo ">> Cleaning..."
	@rm -rf $(BIN_DIR)
	@$(GO) clean
	@echo "🧹 Clean complete"

## Lint the codebase
lint:
	@echo ">> Running linter..."
	@$(LINTER) run $(PKG)

## Run unit tests with race detector and coverage
test:
	@echo ">> Running tests..."
	@$(GO) test -v -race -cover $(PKG)

## Format and tidy modules
tidy:
	@echo ">> Tidying up..."
	@$(GO) fmt $(PKG)
	@$(GO) mod tidy

## Install dependencies
deps:
	@echo ">> Installing dependencies..."
	@$(GO) mod download

## Build Docker image
docker:
	@echo ">> Building Docker image..."
	@docker build \
		--build-arg VERSION=$(GIT_TAG) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		-t $(APP_NAME):$(GIT_TAG) .

## Help menu
help:
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[32m%-12s\033[0m %s\n", $$1, $$2}'
	@echo ""
