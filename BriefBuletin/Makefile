# BriefBuletin Frontend - Angular Makefile

.PHONY: help install build start test watch clean lint

# Default target
help:
	@echo "BriefBuletin Frontend - Angular Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  install    Install npm dependencies"
	@echo "  build      Build production bundle"
	@echo "  start      Start development server"
	@echo "  test       Run unit tests"
	@echo "  watch      Build and watch for changes"
	@echo "  lint       Run linting (if configured)"
	@echo "  clean      Clean build artifacts and dependencies"
	@echo ""

# Install dependencies
install:
	@echo "Installing npm dependencies..."
	npm install

# Build for production
build:
	@echo "Building Angular application for production..."
	npm run build -- --configuration production

# Start development server
start:
	@echo "Starting Angular development server..."
	npm start

# Run tests
test:
	@echo "Running unit tests..."
	npm test

# Build and watch
watch:
	@echo "Building and watching for changes..."
	npm run watch

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@if exist "dist" rmdir /s /q "dist"
	@if exist "node_modules" rmdir /s /q "node_modules"
	@echo "Clean complete."

# Lint (if configured in angular.json)
lint:
	@echo "Running linting..."
	npm run lint 2>/dev/null || echo "Linting not configured or failed"

# Development setup
dev: install
	@echo "Development environment ready. Run 'make start' to begin development."

# Production build
prod: install build
	@echo "Production build complete in dist/ directory"