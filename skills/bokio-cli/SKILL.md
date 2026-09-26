---
name: bokio-cli
description: >
  Interact with the Bokio accounting API via CLI. Manage invoices, customers, items,
  journal entries, credit notes, uploads, bank payments, chart of accounts, fiscal years,
  SIE exports, and connections for Swedish companies. Use when the user wants to perform
  accounting operations, create or list invoices, manage customers, record journal entries,
  or export accounting data from Bokio.
license: MIT
compatibility: Requires Go 1.25+ and the bokio binary (go install github.com/domo84/bokio-cli/cmd/bokio@latest)
metadata:
  author: domo84
  version: "1.0"
  api: https://api.bokio.se/v1/
---

# Bokio CLI Skill

## Overview

`bokio` is a CLI for the Bokio accounting API. It covers all Company API and General API endpoints.

## Prerequisites

1. Install: `go install github.com/domo84/bokio-cli/cmd/bokio@latest`
2. Authenticate with a private token or OAuth:
   ```bash
   bokio auth login --token <TOKEN>
   ```
   Or set `BOKIO_TOKEN` environment variable.
3. Set default company:
   ```bash
   bokio config set company_id <COMPANY_ID>
   ```

## Common workflows

### List and inspect invoices

```bash
bokio invoices list                          # paginated table
bokio invoices list --all -o json            # all invoices as JSON
bokio invoices list -q "status==draft"       # filter drafts
bokio invoices get <id>                      # single invoice
```

### Create and publish an invoice

```bash
bokio invoices create --customer-id <CID> --invoice-date 2025-01-15 --due-date 2025-02-15
bokio invoices line-items add <invoiceId> --description "Consulting" --quantity 10 --unit-price 1500
bokio invoices publish <invoiceId>
```

Or from a JSON file:

```bash
bokio invoices create --from-file invoice.json
```

### Manage customers

```bash
bokio customers list
bokio customers create --name "Acme AB" --email info@acme.se
bokio customers update <id> --from-file customer.json
bokio customers delete <id>
```

### Record a journal entry

```bash
bokio journal-entries create --from-file entry.json
bokio journal-entries reverse <id>
```

Journal entry JSON format:
```json
{
  "date": "2025-01-15",
  "description": "Office supplies",
  "rows": [
    {"accountNumber": 6110, "debitAmount": 1000, "creditAmount": 0},
    {"accountNumber": 1930, "debitAmount": 0, "creditAmount": 1000}
  ]
}
```

### Upload a receipt

```bash
bokio uploads create --file receipt.pdf --description "Office supplies receipt"
```

### Export accounting data

```bash
bokio sie download <fiscalYearId> -O accounting.se
bokio accounts list
bokio fiscal-years list
```

### Record invoice payment

```bash
bokio invoices record <id> --account-number 1920 --payment-date 2025-02-15
```

## Global flags

- `--company-id, -c` — Override company ID
- `--output, -o` — Output format: `table` (default) or `json`
- `--config-dir` — Override config directory

## List flags

All list commands support:
- `--page` — Page number (default: 1)
- `--page-size` — Items per page, max 100 (default: 25)
- `--query, -q` — Filter expression (e.g. `name~John`, `status==draft`, `invoiceDate>=2025-01-01`)
- `--all` — Auto-paginate and fetch everything

## Filter syntax

Operators: `==`, `!=`, `>`, `<`, `>=`, `<=`, `~` (contains). Combine with `&&` (AND) or `||` (OR).

Example: `bokio invoices list -q "status==draft&&invoiceDate>=2025-01-01"`

## Available commands

| Command | Operations |
|---------|-----------|
| `auth` | `login`, `logout`, `status`, `token` |
| `config` | `get`, `set` |
| `company` | get company info |
| `customers` | `list`, `get`, `create`, `update`, `delete` |
| `suppliers` | `list`, `get`, `create`, `update` |
| `items` | `list`, `get`, `create`, `update`, `delete` |
| `invoices` | `list`, `get`, `create`, `update`, `publish`, `record` |
| `invoices line-items` | `add`, `update`, `delete` |
| `invoices attachments` | `add`, `list`, `delete` |
| `invoices payments` | `add`, `list` |
| `invoices settlements` | `add` |
| `journal-entries` | `list`, `get`, `create`, `reverse` |
| `credit-notes` | `list`, `get`, `update`, `publish`, `record` |
| `uploads` | `create`, `list`, `get`, `download` |
| `bank-payments` | `create`, `list`, `get` |
| `accounts` | `list`, `get` |
| `fiscal-years` | `list`, `get` |
| `sie` | `download` |
| `connections` | `list`, `get`, `delete` |

See [references/REFERENCE.md](references/REFERENCE.md) for full command details and flags.
