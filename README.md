# Screen Cammie Chat

Displays the cammie chat along with some other statistics.

## Development Setup

### Prerequisites

1. Go: Check the [.tool-versions](.tool-versions) file for the required Go version.
2. Pre-commit hooks: `git config --local core.hooksPath .githooks/`.
3. Goose (DB migrations): `go install github.com/pressly/goose/v3/cmd/goose@latest`.
4. SQLC (Statically typed queries): `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
5. Make (Optional)

### Configuration

1. Create a `.env` file specifying
  - `APP_ENV`. Available options are:
     -  `development`
     -  `production`
  - `SONG_SPOTIFY_CLIENT_ID`
  - `SONG_SPOTIFY_CLIENT_SECRET`
2. Configure the appropriate settings in the corresponding configuration file located in the [config directory](./config)

## DB

This project uses an SQLite database.
SQLC is used to generate statically typed queries and goose is responsible for the database migrations.

### Usefull commands

- `make migrate`: Run database migrations to update your database schema (watch out, migrations might result in minor data loss).
- `make create-migration`: Create a new migration in the [db/migrations](./db/migrations/) directory.
- `make sqlc`: Generate statically typed queries based on the .sql files in the [db/queries](./db/queries/) directory. Add new queries to this directory as needed.

## Backend

The backend is responsible for fetching and processing external data, which is then stored in the database.
Data can either received by exposing an API or by actively fetching them.

### Running the backend

To build and run the backend, use the following commands:

- `make build-backend`: Build the backend binary.
- `make run-backend`: Run the backend.

### Logs

Backend logs are saved to `./logs/backend.log` (created on first start) and written to `stdout`.

## TUI

The TUI (Text User Interface) displays data retrieved from the database. This flexibility allows for running multiple instances of the TUI, each displaying different data.

### Running the TUI

To build and run the TUI, use the following commands:

- `make build-tui`: Build the TUI binary.
- `make run-tui`: Run the TUI.
-
The TUI requires one argument: the screen name to display. You can create new screens in the [screens directory](./ui/screen/), and you must add them to the startup command list in [tui.go](./internal/cmd/tui.go).

A screen is responsible for creating and managing the layout, consisting of various [views](./ui/view/).

### Logs

TUI logs are written to `./logs/tui.log` and _not_ to `stdout`.
