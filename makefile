all: build

build: clean
	@go build -o scc cmd/tty/scc.go

run:
	@go run cmd/tty/scc.go

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
	@if command -v air > /dev/null; then \
	    air; \
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi
