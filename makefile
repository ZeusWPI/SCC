all: build

build: clean backend tui

run: backend tui
	@./backend & ./tui

backend:
	@[ -f backend ] || (echo "Building backend..." && go build -o backend cmd/backend/backend.go)

tui:
	@[ -f tui ] || (echo "Building tui..." && go build -o tui cmd/tui/tui.go)

clean:
	@rm -f backend tui

sqlc:
	sqlc generate

create-migration:
	@read -p "Enter migration name: " name; \
	goose -dir ./db/migrations create $$name sql

goose:
	@read -p "Action: " action; \
	goose -dir ./db/migrations sqlite3 ./sqlite.db $$action

migrate:
	@goose -dir ./db/migrations sqlite3 ./sqlite.db up

watch:
	@echo "Starting the watch script..."
	./watch.sh
