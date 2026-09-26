# bokio-cli

A command-line interface for the [Bokio](https://www.bokio.se) accounting API.

Manage invoices, customers, journal entries, and more from your terminal.

## Install

```bash
go install github.com/domo84/bokio-cli/cmd/bokio@latest
```

Or build from source:

```bash
git clone git@github.com:domo84/bokio-cli.git
cd bokio-cli
go build -o bokio ./cmd/bokio/
```

## Authentication

### Private token

Generate a token in the Bokio app, then:

```bash
bokio auth login --token <your-token>
```

Or set the environment variable:

```bash
export BOKIO_TOKEN=<your-token>
```

### OAuth 2.0

For public integrations using OAuth:

```bash
bokio config set client_id <your-client-id>
bokio config set client_secret <your-client-secret>
bokio auth login --oauth
```

### Check status

```bash
bokio auth status
```

## Configuration

Set your default company ID to avoid passing `-c` on every command:

```bash
bokio config set company_id <your-company-id>
```

Config is stored at `~/.config/bokio-cli/config.yaml`. Use `--config-dir` to keep it
elsewhere.

Available settings:

| Key | Description | Default | Env var |
|-----|-------------|---------|---------|
| `company_id` | Default company ID, so `-c` can be omitted | | `BOKIO_COMPANY_ID` |
| `output_format` | Output format (`table` or `json`) | `table` | `BOKIO_OUTPUT` |
| `client_id` | OAuth client ID | | `BOKIO_CLIENT_ID` |
| `client_secret` | OAuth client secret | | `BOKIO_CLIENT_SECRET` |
| `redirect_port` | Local port for the OAuth callback | `8585` | `BOKIO_REDIRECT_PORT` |

Settings resolve in this order, highest first: **command-line flags**, **`BOKIO_*`
environment variables**, **the config file**, **built-in defaults**.

To see every key with its current value and where that value came from:

```bash
bokio config list
```

`bokio config --help` documents the same list, so it is always available offline.

## Usage

### Company info

```bash
bokio company
```

### Customers

```bash
bokio customers list
bokio customers get <id>
bokio customers create --name "Acme AB" --email info@acme.se --city Stockholm
bokio customers update <id> --from-file customer.json
bokio customers delete <id>
```

### Suppliers

```bash
bokio suppliers list
bokio suppliers list -q "name==Supplier ABC"
bokio suppliers get <id>
bokio suppliers update <id> --bankgiro 5097-1282
bokio suppliers update <id> --from-file supplier.json
```

### Invoices

```bash
bokio invoices list
bokio invoices list --all                    # auto-paginate
bokio invoices list -q "status==draft"       # filter
bokio invoices get <id>
bokio invoices create --customer-id <id> --invoice-date 2025-01-15 --due-date 2025-02-15
bokio invoices create --from-file invoice.json
bokio invoices publish <id>
bokio invoices record <id> --account-number 1920 --payment-date 2025-02-15
```

Invoice sub-resources:

```bash
bokio invoices line-items add <invoiceId> --description "Consulting" --quantity 10 --unit-price 1500
bokio invoices attachments list <invoiceId>
bokio invoices payments list <invoiceId>
bokio invoices payments add <invoiceId> --amount 15000 --payment-date 2025-02-15 --account-number 1920
bokio invoices settlements add <invoiceId> --amount 15000 --settle-date 2025-02-15 --account-number 1920
```

### Items

```bash
bokio items list
bokio items create --description "Consulting hour" --unit-price 1500 --unit hours --vat-rate 25
bokio items update <id> --from-file item.json
bokio items delete <id>
```

### Journal entries

```bash
bokio journal-entries list
bokio journal-entries get <id>
bokio journal-entries create --from-file entry.json
bokio journal-entries reverse <id>
```

### Credit notes

```bash
bokio credit-notes list
bokio credit-notes get <id>
bokio credit-notes publish <id>
bokio credit-notes record <id> --account-number 1920 --payment-date 2025-03-01
```

### Uploads

```bash
bokio uploads list
bokio uploads create --file receipt.pdf --description "Office supplies"
bokio uploads download <id> -O receipt.pdf
```

### Bank payments

```bash
bokio bank-payments list
bokio bank-payments get <id>
bokio bank-payments create --from-file payment.json
```

### Chart of accounts

```bash
bokio accounts list
bokio accounts get 1920
```

### Fiscal years

```bash
bokio fiscal-years list
```

### SIE export

```bash
bokio sie download <fiscalYearId> -O accounting.se
```

### Connections

```bash
bokio connections list
bokio connections delete <id>
```

## Global flags

| Flag | Short | Description |
|------|-------|-------------|
| `--company-id` | `-c` | Company ID (or set `company_id` in config / `BOKIO_COMPANY_ID`) |
| `--output` | `-o` | Output format: `table` or `json` (or set `output_format` in config) |
| `--config-dir` | | Config directory path (default `~/.config/bokio-cli`) |

List commands also support:

| Flag | Short | Description |
|------|-------|-------------|
| `--page` | | Page number (default: 1) |
| `--page-size` | | Items per page, max 100 (default: 25) |
| `--query` | `-q` | Filter expression (e.g. `name~John`, `status==draft`) |
| `--all` | | Fetch all pages automatically |

## Shell completion

```bash
# Bash
source <(bokio completion bash)

# Zsh
source <(bokio completion zsh)

# Fish
bokio completion fish | source
```

## API reference

https://docs.bokio.se/reference/overview

## License

MIT
