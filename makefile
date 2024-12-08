# Variables
BACKEND_BIN := backend_bin
TUI_BIN := tui_bin
BACKEND_SRC := cmd/backend/backend.go
TUI_SRC := cmd/tui/tui.go
DB_DIR := ./db/migrations
DB_USER := postgres
DB_PASSWORD := postgres
DB_NAME := scc

# Phony targets
.PHONY: all build build-backed build-tui clean run run-backend run-tui db sqlc create-migration goose migrate watch

# Default target: build everything
all: build

# Build targets
build: clean build-backend build-tui

build-backend: clean-backend
	@echo "Building $(BACKEND_BIN)..."
	@go build -o $(BACKEND_BIN) $(BACKEND_SRC)

build-tui: clean-tui
	@echo "Building $(TUI_BIN)..."
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

# Run db
db:
	@docker compose up

# Clean targets
clean: clean-backend clean-tui

clean-backend:
	@if [ -f "$(BACKEND_BIN)" ]; then \
		echo "Cleaning $(BACKEND_BIN)..."; \
		rm -f "$(BACKEND_BIN)"; \
	fi

clean-tui:
	@if [ -f "$(TUI_BIN)" ]; then \
		echo "Cleaning $(TUI_BIN)..."; \
		rm -f "$(TUI_BIN)"; \
	fi

# SQL and migration targets
sqlc:
	sqlc generate

create-migration:
	@read -p "Enter migration name: " name; \
	goose -dir $(DB_DIR) create $$name sql

migrate:
	@goose -dir $(DB_DIR) postgres "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=localhost sslmode=disable" up

goose:
	@read -p "Action: " action; \
	goose -dir $(DB_DIR) postgres "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=localhost sslmode=disable" $$action
