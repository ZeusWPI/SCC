# Variables
BACKEND_BIN := backend
TUI_BIN := tui
BACKEND_SRC := cmd/backend/backend.go
TUI_SRC := cmd/tui/tui.go
DB_DIR := ./db/migrations
DB_FILE := ./sqlite.db

# Phony targets
.PHONY: all build clean run run-backend run-tui sqlc create-migration goose migrate watch

# Default target: build everything
all: build

# Build targets
build: clean build-backend build-tui

build-backend:
	@echo "Building $(BACKEND_BIN)..."
	@rm -f $(BACKEND_BIN)
	@go build -o $(BACKEND_BIN) $(BACKEND_SRC)

build-tui:
	@echo "Building $(TUI_BIN)..."
	@rm -f $(TUI_BIN)
	@go build -o $(TUI_BIN) $(TUI_SRC)

# Run targets
run: run-backend run-tui

run-backend:
	@[ -f $(BACKEND_BIN) ] || $(MAKE) build-backend
	@./$(BACKEND_BIN)

run-tui:
	@[ -f $(TUI_BIN) ] || $(MAKE) build-tui
	@read -p "Enter screen name: " screen; \
	./$(TUI_BIN) $$screen

# Clean targets
clean: clean-backend clean-tui

clean-backend:
	@echo "Cleaning $(BACKEND_BIN)..."
	@rm -f $(BACKEND_BIN)

clean-tui:
	@echo "Cleaning $(TUI_BIN)..."
	@rm -f $(TUI_BIN)

# SQL and migration targets
sqlc:
	sqlc generate

create-migration:
	@read -p "Enter migration name: " name; \
	goose -dir $(DB_DIR) create $$name sql

goose:
	@read -p "Action: " action; \
	goose -dir $(DB_DIR) sqlite3 $(DB_FILE) $$action

migrate:
	@goose -dir $(DB_DIR) sqlite3 $(DB_FILE) up
