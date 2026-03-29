# AGENTS.md

Context for AI coding agents working on this codebase.

## Project overview

A Go CLI that wraps the Bokio accounting API (`https://api.bokio.se/v1/`). Built with cobra for commands, viper for config. Module path: `github.com/domo84/bokio-cli`.

## Setup commands

```bash
go mod tidy
go build -o bokio ./cmd/bokio/
```

## Build and test

```bash
go build ./...           # compile all packages
go vet ./...             # static analysis
go build -o bokio ./cmd/bokio/   # produce binary
```

No test suite yet. When adding tests, use `go test ./...`.

## Project structure

```
cmd/bokio/main.go          Entry point, calls commands.Execute()
commands/                   Cobra command definitions (one file per resource)
  root.go                   Root command, global flags, initState/initStateWithClient
  helpers.go                Shared list flags (page, page-size, query, all)
  auth.go                   auth login|logout|status|token
  config.go                 config get|set
  company.go                company (get info)
  customers.go              customers list|get|create|update|delete
  items.go                  items list|get|create|update|delete
  invoices.go               invoices + sub-resources (line-items, attachments, payments, settlements)
  journal_entries.go        journal-entries list|get|create|reverse
  credit_notes.go           credit-notes list|get|update|publish|record
  uploads.go                uploads create|list|get|download
  bank_payments.go          bank-payments create|list|get
  chart_of_accounts.go      accounts list|get
  fiscal_years.go           fiscal-years list|get
  sie.go                    sie download
  connections.go            connections list|get|delete
  completion.go             Shell completion (bash/zsh/fish/powershell)
internal/
  api/
    client.go               HTTP client (Get, Post, Put, Delete + JSON helpers, rate limit retry)
    models.go               All request/response structs
    pagination.go           PaginatedResponse[T], AutoPaginate, ListParams
    errors.go               BokioError type
    <resource>.go           One file per API resource with methods on *Client
  auth/
    auth.go                 OAuth 2.0 flow (PKCE, client credentials, refresh)
    token_store.go          File-based credential storage (~/.config/bokio-cli/credentials.json)
    server.go               Local HTTP callback server for OAuth redirect
  config/
    config.go               Viper-based config (~/.config/bokio-cli/config.yaml)
  output/
    formatter.go            Formatter interface, NewFormatter factory
    json.go                 JSON output
    table.go                Table output with type-switch rendering per model
```

## Code style

- Go standard formatting (`gofmt`)
- No unnecessary abstractions; each resource is a flat file
- Pointer receivers on `*Client` for API methods
- Use `any` not `interface{}`
- Errors returned, not panicked
- cobra commands constructed via `newXxxCmd()` factory functions

## Key patterns

### Adding a new command

1. Create `internal/api/<resource>.go` with methods on `*Client`
2. Create `commands/<resource>.go` with cobra commands
3. Register in `commands/root.go` `Execute()` via `root.AddCommand(newXxxCmd())`
4. Add table rendering case in `internal/output/table.go`

### Command structure

Every command that calls the API follows this pattern:

```go
state, err := initStateWithClient(cmd)   // loads config, resolves token, creates API client
if err != nil { return err }
if err := requireCompanyID(state); err != nil { return err }  // for company-scoped endpoints
result, err := state.client.SomeMethod(cmd.Context(), ...)
if err != nil { return err }
return state.formatter.Format(result)
```

Use `initState(cmd)` instead of `initStateWithClient(cmd)` for commands that don't need the API client (auth, config).

### API client

- `Client.companyURL(path)` builds `/companies/{companyId}/path`
- `Client.generalURL(path)` builds `/path` (for auth, connections)
- JSON helpers: `GetJSON`, `PostJSON`, `PutJSON`, `PostEmpty`, `DeleteEmpty`
- Multipart: `PostRaw` (used by uploads)
- Rate limit: auto-retries on 429 using `Bokio-RateLimit-RetryAfter` header

### Pagination

API returns `{ items: [], totalItems, totalPages, currentPage }`.
`AutoPaginate[T]()` generic fetches all pages.
List commands support `--all` flag to auto-paginate or `--page`/`--page-size` for manual.

### Config and auth

- Config: `~/.config/bokio-cli/config.yaml` (viper, `BOKIO_*` env overrides)
- Credentials: `~/.config/bokio-cli/credentials.json` (0600 permissions)
- Token resolution order: `BOKIO_TOKEN` env var > stored credentials

## Security considerations

- Never log or print tokens except in `bokio auth token` (explicit user request)
- Credentials file must use 0600 permissions, config dir 0700
- Always use HTTPS (hardcoded base URL)
- Do not commit `.env` or `credentials.json` (covered by `.gitignore`)

## Commit guidelines

- Default branch: `trunk`
- Conventional commit style preferred
- Co-author AI contributions

## API reference

https://docs.bokio.se/reference/overview
