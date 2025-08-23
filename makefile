# Variables
BACKEND_BIN := backend_bin
TUI_BIN := tui_bin
BACKEND_SRC := cmd/backend/backend.go
TUI_SRC := cmd/tui/tui.go
DB_DIR := ./db/migrations
DB_USER := postgres
DB_PASSWORD := postgres
DB_NAME := scc

all: build

setup:
	@go get tool

build: clean build-backend build-tui

build-backend: clean-backend
	@echo "Building $(BACKEND_BIN)..."
	@go build -o $(BACKEND_BIN) $(BACKEND_SRC)

build-tui: clean-tui
	@echo "Building $(TUI_BIN)..."
	@go build -o $(TUI_BIN) $(TUI_SRC)

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

backend:
	@docker compose up -d
	@[ -f $(BACKEND_BIN) ] || $(MAKE) build-backend
	@./$(BACKEND_BIN)
	@docker compose down

tui:
	@[ -f $(TUI_BIN) ] || $(MAKE) build-tui
	@read -p "Enter screen name: " screen; \
	./$(TUI_BIN) $$screen

goose:
	@docker compose down
	@docker compose up db -d
	@docker compose exec db bash -c 'until pg_isready -U postgres; do sleep 1; done'
	@read -p "Action: " action; \
	go tool goose -dir ./db/migrations postgres "user=postgres password=postgres host=localhost port=5431 dbname=website sslmode=disable" $$action
	@docker compose down db

migrate:
	@docker compose down
	@docker compose up db -d
	@docker compose exec db bash -c 'until pg_isready -U postgres; do sleep 1; done'
	@go tool goose -dir ./db/migrations postgres "user=postgres password=postgres host=localhost port=5431 dbname=website sslmode=disable" up
	@docker compose down db

create-migration:
	@read -p "Enter migration name: " name; \
	go tool goose -dir ./db/migrations create $$name sql

query:
	@go tool sqlc generate

dead:
	@go tool deadcode ./...

.PHONY: all setup build build-backed build-tui clean clean-backend clean-tui backend tui create-migration goose query dead
