# Gator

A CLI tool for aggregating RSS feeds.

## Requirements

- Go (version 1.26 or later)
- PostgreSQL

## Installation

Install the gator CLI using go install:

```bash
go install github.com/carlbark/gator@latest
```

## Setup

1. Ensure PostgreSQL is running and create a database named `gator`.

2. Install goose (database migration tool):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

3. Run the database migrations from the `sql/schema` directory:

```bash
cd sql/schema
goose postgres "postgres://your_username:your_password@localhost:5432/gator?sslmode=disable" up
```

Replace `your_username` and `your_password` with your actual PostgreSQL credentials.

4. Create a config file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "your_username"
}
```

Replace `username`, `password`, and `your_username` with your actual values.

## Usage

After setup, you can run various commands:

- `gator register <username>` - Register a new user
- `gator login <username>` - Login as a user
- `gator addfeed <name> <url>` - Add a new RSS feed
- `gator feeds` - List all feeds
- `gator agg <time>` - Aggregate feeds for a certain time (e.g., 1s, 1m)- This runs as an infinite loop (CTRL-C to exit)
- `gator browse <limit>` - Browse posts (default limit 2)

For more commands, run `gator` without arguments to see available commands.