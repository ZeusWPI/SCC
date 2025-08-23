# Screen Cammie Chat

A terminal based dashboard for viewing the Cammie chat messages and other Zeus related data.

## Overview

The project has 2 main parts:

### 1. TUI

A terminal interface built with [bubbletea](https://github.com/charmbracelet/bubbletea).
It displays the Cammie chat messages and other Zeus data in a screen-based layout.

- **Views**: Reusable components responsible for rendering a single type of data (e.g. the [TAP](github.com/zeusWPI/tap) statistics).
- **Screens**: A combination of multiple views forming a single terminal screen.

Each view implements a shared interface, exposing methods for initialization, updating, and rendering.
When a screen is loaded:

- Each view’s initial data is fetched automatically.
- Each view’s update loop runs in a goroutine, periodically checking for new data.
- The update loop communicates state changes back to the TUI.

This design allows each view to be self-contained and independently refresh its data.

However, not all data can be retrieved directly from persistent external services and that’s where the backend comes in.

### 2. Backend

The backend handles data aggregation and persistence.

- Exposes an API for incoming data.
- Maintaines websockets.
- Makes use of external API's.

Provides persistence data using Postgres for data that the external services do not retain.

## Development

### Prerequisites

1. Download the golang version [.tool-versions](.tool-versions).
2. Install make.
3. Install the go tools: `make setup`.
4. (Optional) Install the pre-commit hooks: `git config --local core.hooksPath .githooks/`.

### Configuration

1. Copy `.env.example` to `.env`, set `ENV=development` and populate the remaining keys.
2. (Optional) Edit the [config file](./config/development.yml). The defaults work.

### Database

A Postgres database instance is provided via docker compose and started automatically when needed by the [makefile](./makefile).
To use a custom database, update the config and edit the makefile.

### Run

1. Migrate the database `make migrate`.
2. Start the backend `make backend`.
3. Start a TUI `make tui` and enter the desired screen name (if you're not using the makefile use the `-screen` flag to specify the screen).

### Logs

- Backend: written to `./logs/backend.log` and to stdout.
- TUI: only written to `./logs/{screen}.log`.

### Useful commands

- `make create-migration`: Create a new migration in the [db/migrations](./db/migrations/) directory.
- `make query`: Generate statically typed queries based on the .sql files in the [db/queries](./db/queries/) directory. Add new queries to this directory as needed.
- `make goose`: Migrate one version up or down.
- `make dead`: Check for unreachable code.

## Production

1. Set `ENV=production` in `.env`.
2. Provide a Postgres database.
3. Populate the [production config file](./config/production.yml).
4. Build the binaries `make build`.
5. Run both binaries with the desired flags.
