# BriefBuletin - Combined Makefile

.PHONY: help install-all build-all run-all clean

# Default target
help:
	@echo "BriefBuletin - Complete News Platform Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  install-all    Install dependencies for all components"
	@echo "  build-all      Build all components"
	@echo "  run-all        Run all services (API, Frontend, Scraper)"
	@echo "  clean          Clean build artifacts"
	@echo ""
	@echo "Individual component targets:"
	@echo "  install-frontend    Install Angular dependencies"
	@echo "  build-frontend      Build Angular application"
	@echo "  run-frontend        Start Angular dev server"
	@echo "  install-api         Download Go modules"
	@echo "  build-api           Build Go API binary"
	@echo "  run-api             Start Go API server"
	@echo "  install-scraper     Install Python dependencies"
	@echo "  run-scraper         Start news scraper (background)"
	@echo "  fetch-scraper       Run scraper once"
	@echo ""

# Install all dependencies
install-all: install-frontend install-api install-scraper

# Build all components
build-all: build-frontend build-api

# Run all services
run-all:
	@echo "Starting all services..."
	@echo "Note: Services will run in foreground. Use Ctrl+C to stop all."
	@start /B make run-api &
	@timeout /t 5 /nobreak > nul
	@start /B make run-frontend &
	@timeout /t 5 /nobreak > nul
	@make run-scraper-background
	@echo "All services started. Press Ctrl+C to stop."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@if exist "BriefBuletin\dist" rmdir /s /q "BriefBuletin\dist"
	@if exist "BriefBuletin_API\build" rmdir /s /q "BriefBuletin_API\build"
	@if exist "BriefBuletin_NewsScraper\build" rmdir /s /q "BriefBuletin_NewsScraper\build"
	@echo "Clean complete."

# Frontend targets
install-frontend:
	@echo "Installing Angular dependencies..."
	cd BriefBuletin && npm install

build-frontend:
	@echo "Building Angular application..."
	cd BriefBuletin && npm run build -- --configuration production

run-frontend:
	@echo "Starting Angular dev server..."
	cd BriefBuletin && npm start

# API targets
install-api:
	@echo "Downloading Go modules..."
	cd BriefBuletin_API && go mod download

build-api:
	@echo "Building Go API..."
	cd BriefBuletin_API && make winBuild

run-api:
	@echo "Starting Go API server..."
	cd BriefBuletin_API && make dev

# Scraper targets
install-scraper:
	@echo "Installing Python dependencies..."
	cd BriefBuletin_NewsScraper && pip install feedparser requests beautifulsoup4 psycopg2-binary transformers torch python-dateutil

run-scraper:
	@echo "Starting news scraper (background)..."
	cd BriefBuletin_NewsScraper && make start

run-scraper-background:
	@echo "Starting news scraper in background..."
	cd BriefBuletin_NewsScraper && make start

fetch-scraper:
	@echo "Running scraper once..."
	cd BriefBuletin_NewsScraper && make fetch

# Development helpers
dev-setup: install-all
	@echo "Development environment setup complete."
	@echo "Run 'make run-all' to start all services."

# Status check
status:
	@echo "Checking service status..."
	@tasklist /FI "IMAGENAME eq briefbuletin.exe" /NH >nul 2>&1 && echo "API: Running" || echo "API: Not running"
	@tasklist /FI "IMAGENAME eq node.exe" /NH >nul 2>&1 && echo "Frontend: Running" || echo "Frontend: Not running"
	@cd BriefBuletin_NewsScraper && make status