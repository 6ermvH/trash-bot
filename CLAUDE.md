# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Run the application
go run ./cmd/main.go

# Build
go build -v ./...

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -v ./internal/repository/inmemory -run TestFunctionName

# Lint (requires golangci-lint v2)
golangci-lint run

# Format code (requires gofumpt, goimports, golines)
gofumpt -w .
goimports -w .
golines -w --max-len=120 .
```

## Architecture Overview

This is a Telegram bot for managing trash duty rotation with an optional HTTP admin panel.

### Layer Structure

```
cmd/main.go              → Entry point, initializes components and runs them via errgroup
    ↓
cmd/bot/bot.go           → Telegram bot setup and registration of handlers
cmd/panel/panel.go       → HTTP server (Gin) with embedded web assets
    ↓
internal/services/trashmanager/  → Business logic (Who/Next/Prev/SetEstablish)
    ↓
internal/repository/     → Data access layer (interface + implementations)
    ├── inmemory/        → Map-based storage (loses data on restart)
    └── sqlite/          → Persistent SQLite storage
```

### Key Components

- **TrashManager Service** (`internal/services/trashmanager/`): Core business logic for managing user rotation per chat. Methods: `Who`, `Next`, `Prev`, `SetEstablish`, `Stats`, `Chats`, `Subscribe`, `Unsubscribe`.

- **Scheduler Service** (`internal/services/scheduler/`): Background service that sends daily reminders to subscribed chats at their configured time. Runs on a per-minute ticker.

- **Telegram Handlers** (`internal/handlers/telegram/`): Bot commands (`/start`, `/set`, `/next`, `/prev`, `/who`, `/subscribe`, `/unsubscribe`) and inline keyboard callbacks.

- **HTTP Handlers** (`internal/handlers/http/v1/`): REST API with JWT authentication. Endpoints: `/api/login`, `/api/stats`, `/api/chats`, `/api/chats/:id`.

- **Repository Pattern**: `Chat` model with `ID` (Telegram chat ID), `Current` (index), `Users` (list of names), `NotifyTime` (optional, for daily reminders). Storage backend selected via `config.Database.Type`.

### Configuration

Edit `config/base.yaml`:
- `telegram.botkey`: Telegram bot token
- `server.enabled`: Enable/disable HTTP panel
- `database.type`: `"sqlite"` for persistent storage, anything else for in-memory

### Concurrency

Bot and HTTP panel run concurrently using `errgroup`. Both receive context from main for graceful shutdown.
