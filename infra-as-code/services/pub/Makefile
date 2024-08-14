APP_NAME := publisher
BUILD_DIR := bin

GO := go
GO_BUILD := $(GO) build
GO_RUN := $(GO) run
GO_TEST := $(GO) test
GO_FMT := $(GO) fmt
GO_LINT := golangci-lint run
GO_CLEAN := $(GO) clean
GO_MOD_TIDY := $(GO) mod tidy

# Default target
all: build

# Build the project
build:
	@echo "Building $(APP_NAME)..."
	$(GO_BUILD) -o $(BUILD_DIR)/$(APP_NAME) .

container:
	@echo "Building container image..."
	docker build --platform linux/amd64 -t $(APP_NAME) .

# Run the project
run:
	@echo "Running $(APP_NAME)..."
	$(GO_RUN) main.go

# Test the project
test:
	@echo "Running tests..."
	$(GO_TEST) ./...

# Format the code
fmt:
	@echo "Formatting code..."
	$(GO_FMT) ./...

# Lint the code (requires golangci-lint to be installed)
lint:
	@echo "Linting code..."
	$(GO_LINT)

# Clean the build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GO_CLEAN)
	rm -rf $(BUILD_DIR)

# Tidy up the module
tidy:
	@echo "Tidying up go.mod..."
	$(GO_MOD_TIDY)