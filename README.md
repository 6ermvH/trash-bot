# Trash Bot

[![CI](https://img.shields.io/github/actions/workflow/status/6ermvH/trash-bot/go.yml?branch=main)](https://github.com/6ermvH/trash-bot/actions)
[![Tag](https://img.shields.io/github/v/tag/6ermvH/trash-bot)](https://github.com/6ermvH/trash-bot/tags)
[![Go Version](https://img.shields.io/github/go-mod/go-version/6ermvH/trash-bot)](https://github.com/6ermvH/trash-bot/blob/main/go.mod)
[![License](https://img.shields.io/github/license/6ermvH/trash-bot)](https://github.com/6ermvH/trash-bot/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/6ermvH/trash-bot)](https://goreportcard.com/report/github.com/6ermvH/trash-bot)

Telegram bot for managing a trash duty rotation, with an optional admin panel.

## Features
- Telegram commands: `/start`, `/set`, `/next`, `/prev`, `/who`, `/subscribe`, `/unsubscribe`
- Daily notifications at a user-selected time
- SQLite or in-memory storage for chat state
- Optional HTTP admin panel (Gin) with JWT authentication

## Requirements
- Go (see `go.mod` for the exact version)
- Optional tooling for formatting and linting:
  - `golangci-lint`
  - `golines`, `gofumpt`, `goimports`

## Configuration
Edit `config/base.yaml` with your bot token and server settings.

```yaml
telegram:
  botkey: "<your-telegram-bot-token>"

server:
  enabled: true
  addr: "127.0.0.1"
  port: "8080"
  adminlogin: "admin"
  adminpassword: "admin"
  jwtsecret: "your-secret-key"

database:
  type: "sqlite"  # or leave empty for in-memory
  path: "data/trash.db"
```

## Run
```bash
go run ./cmd/main.go

```
## Tests
```bash
go test ./...
```

## Project structure
- `cmd/` - application entry points
- `internal/` - core business logic, handlers, services, repositories
- `config/` - configuration files
