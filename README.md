# trash-bot

[![CI](https://img.shields.io/github/actions/workflow/status/6ermvH/trash-bot/ci.yml?branch=main)](https://github.com/6ermvH/trash-bot/actions)
[![Release](https://img.shields.io/github/v/release/6ermvH/trash-bot)](https://github.com/6ermvH/trash-bot/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/6ermvH/trash-bot)](https://github.com/6ermvH/trash-bot/blob/main/go.mod)
[![License](https://img.shields.io/github/license/6ermvH/trash-bot)](https://github.com/6ermvH/trash-bot/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/6ermvH/trash-bot)](https://goreportcard.com/report/github.com/6ermvH/trash-bot)
[![Codecov](https://codecov.io/gh/6ermvH/trash-bot/branch/main/graph/badge.svg)](https://codecov.io/gh/6ermvH/trash-bot)

Telegram bot for managing a trash duty rotation, with an optional admin panel.

## Features
- Telegram commands: `/start`, `/set`, `/next`, `/prev`, `/who`
- In-memory storage for chat state
- Optional HTTP panel (Gin) for future CRM/admin features

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
