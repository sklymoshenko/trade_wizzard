# Variables
APP_NAME := trade_wizzard
BUILD_DIR := ./bin

# Default target
.PHONY: all
all: build

# Build the project
.PHONY: build
build:
	@echo "Building the project..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

# Run the project
.PHONY: run
run: build
	@echo "Running the project..."
	@$(BUILD_DIR)/$(APP_NAME)

# Clean up build files
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Format code
.PHONY: format
format:
	@echo "Formatting code..."
	@go fmt ./...