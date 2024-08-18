# Variables
APP_NAME = repocate
SRC_DIR = ./cmd/repocate
BUILD_DIR = ./build

# Targets
.PHONY: all clean build install test

all: clean build install

clean:
	rm -rf $(BUILD_DIR)

build:
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

install:
	@echo "Installing $(APP_NAME)..."
	mv $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/

test:
	@echo "Running tests..."
	go test ./...

docker:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

run:
	@echo "Running $(APP_NAME)..."
	$(BUILD_DIR)/$(APP_NAME)