# Screen Cammie Chat

Displays the cammie chat along with some other statistics.

## Development

Check [.tool-versions](.tool-versions) for the current used version of golang

- Install pre-commit hooks `git config --local core.hooksPath .githooks/`.
- Install goose `go install github.com/pressly/goose/v3/cmd/goose@latest`.
- Install sqlc `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- Install air `go install github.com/air-verse/air@latest`

Create a `.env` (Look at [.env.example](.env.example])).

Start developing with `make build` and `make run`.
For live reloading `inotify` is required `sudo apt install inotify-tools`.
`make watch` will start the hot reloading.

Logs will be logged to the `./logs` directory (will be made at the start of the first run).

## Build & Run

- `make`
- `make run`
