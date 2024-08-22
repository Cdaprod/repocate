# Project variables
PROJECT_NAME = repocate
BUILD_DIR = ./build
BIN_DIR = ./bin
SRC_DIR = ./cmd/main
INSTALL_DIR = /usr/local/bin
MAN_DIR = /usr/local/share/man/man1
MAN_PAGE = repocate.1
REPO_DIR = /opt/repocate
PROJECT_DIR = $(HOME)/Repocate

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

# Install the binary, man page, and set up the repository directory
install: build
	@echo "Installing $(PROJECT_NAME)..."
	@sudo install -m 0755 $(BIN_DIR)/$(PROJECT_NAME) $(INSTALL_DIR)
	@echo "Installing man page..."
	@sudo install -Dm644 $(MAN_PAGE) $(MAN_DIR)/$(MAN_PAGE)
	@echo "Setting up repository directory..."
	@sudo mkdir -p $(REPO_DIR)
	@sudo cp -r . $(REPO_DIR) # Copy all files to the repository directory
	@sudo chown -R $(USER):$(USER) $(REPO_DIR) # Change ownership to the current user
	@echo "Repository directory setup complete."
	@echo "Setting up project directory..."
	@mkdir -p $(PROJECT_DIR) # Create the project directory in the user's home
	@chown -R $(USER):$(USER) $(PROJECT_DIR) # Change ownership to the current user
	@echo "Project directory setup complete."
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