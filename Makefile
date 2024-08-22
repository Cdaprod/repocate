# Project variables
PROJECT_NAME = repocate
BUILD_DIR = ./build
BIN_DIR = ./bin
SRC_DIR = ./cmd/main
INSTALL_DIR = /usr/local/bin

# Go-related variables
GO = go
GO_FLAGS = -v
BUILD_FLAGS = -o
TEST_FLAGS = ./...
LINT_FLAGS = ./...

# Default target: Build the project
all: clean build

# Clean up the build directory
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR) $(BUILD_DIR)
	@echo "Clean complete."

# Build the project
build: clean
	@echo "Building $(PROJECT_NAME)..."
	@$(GO) build $(GO_FLAGS) $(BUILD_FLAGS) $(BIN_DIR)/$(PROJECT_NAME) $(SRC_DIR)
	@echo "Build complete."

# Install the binary to /usr/local/bin
install: build
	@echo "Installing $(PROJECT_NAME)..."
	@sudo install -m 0755 $(BIN_DIR)/$(PROJECT_NAME) $(INSTALL_DIR)
	@echo "Install complete."

# Test the project
test:
	@echo "Running tests..."
	@$(GO) test $(TEST_FLAGS)
	@echo "Tests complete."

# Run linting
lint:
	@echo "Linting code..."
	@$(GO) vet $(LINT_FLAGS)
	@echo "Lint complete."

# Run the binary
run: build
	@echo "Running $(PROJECT_NAME)..."
	@./$(BIN_DIR)/$(PROJECT_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@$(GO) mod tidy
	@echo "Dependencies installed."

# Default target to build and test
.PHONY: all clean build install test lint run deps