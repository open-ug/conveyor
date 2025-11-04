# Conveyor Makefile -

# Version constants
SWAG_VERSION := v1.16.6
# Default version for builds, can be overridden (e.g., make deb VERSION=1.0.0).
# The 'export' keyword makes it available to sub-processes like dpkg-buildpackage.
export VERSION ?= snapshot

# Ensure that targets which are not files are always executed.
.PHONY: help install-swag swagger-init start test docs deb clean-deb clean

# Default target
help: ## Show this help message
	@echo "Conveyor - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Development Targets ---

start: ## Start the API server
	@echo "Starting API server..."
	go run main.go up

test: ## Run tests
	@echo "Running tests..."
	APP_ENV=test go test ./... -v

docs: ## Start documentation server
	@echo "Starting documentation server..."
	cd docs && npm run start

# --- Tooling & Documentation Generation ---

install-swag: ## Install swag CLI tool
	@echo "Installing swag CLI tool version $(SWAG_VERSION)..."
	go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
	@echo "swag CLI tool installed"

swagger-init: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/app.go -o internal/swagger
	@echo "Swagger documentation generated in internal/swagger/"

# --- Packaging Targets ---

deb: clean-deb ## Build Debian package
	@echo "Building Debian package (VERSION=$(VERSION))..."
	@# Copy packaging files to the required 'debian' directory at the root.
	cp -R packaging/debian debian
	@# Build the binary-only, unsigned package. The exported VERSION var is used by debian/rules.
	dpkg-buildpackage -us -uc -b
	@echo "Debian package built."
	@# Clean up the temporary directory used for the build.
	rm -rf debian

# --- Clean Targets ---

clean-deb: ## Clean Debian build artifacts
	@echo "Cleaning Debian build artifacts..."
	@# Remove the temporary 'debian' directory.
	rm -rf debian
	@# dpkg-buildpackage places artifacts in the parent directory of the source tree.
	@# We remove them explicitly. Using a package-specific name is safer than a generic glob.
	@# The '|| true' prevents errors if no files are found.
	rm -f ../conveyor_*.deb ../conveyor_*.ddeb ../conveyor_*.changes ../conveyor_*.buildinfo ../conveyor_*.dsc || true
	@echo "Debian build artifacts cleaned."

clean: clean-deb ## Clean all project build artifacts
	@echo "Cleaning all build artifacts..."
	@# Add other project-specific clean rules here if they exist.