all: build

build: clean
	@go build -o scc cmd/tty/scc.go

run:
	@./scc

clean:
	@rm -f scc

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
