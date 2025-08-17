# Conveyor Makefile -

# Version constants
SWAG_VERSION := v1.16.6

.PHONY: help install-swag swagger-init start

# Default target
help: ## Show this help message
	@echo "Conveyor - Swagger Documentation Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Swagger documentation commands
swagger-init: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/app.go -o internal/swagger
	@echo "Swagger documentation generated in docs/swagger/"

install-swag: ## Install swag CLI tool
	@echo "Installing swag CLI tool version $(SWAG_VERSION)..."
	go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
	@echo "swag CLI tool installed" 

start: ## Start the API server
	@echo "Starting API server..."
	go run main.go api-server

test: ## Run tests
	@echo "Running tests..."
	go test ./... -v