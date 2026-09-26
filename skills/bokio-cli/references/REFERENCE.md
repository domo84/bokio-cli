# Bokio CLI Command Reference

## Auth

```
bokio auth login --token <TOKEN>       # Private token auth
bokio auth login --oauth               # OAuth 2.0 flow (requires client_id/client_secret in config)
bokio auth logout                      # Remove stored credentials
bokio auth status                      # Show auth state
bokio auth token                       # Print current access token (for scripting)
```

## Config

```
bokio config get <key>                 # Read a config value
bokio config set <key> <value>         # Write a config value
```

Valid keys: `auth_mode`, `company_id`, `output_format`, `client_id`, `client_secret`, `redirect_port`

Config file: `~/.config/bokio-cli/config.yaml`

Environment overrides: `BOKIO_TOKEN`, `BOKIO_COMPANY_ID`, `BOKIO_OUTPUT`, `BOKIO_AUTH_MODE`, `BOKIO_CLIENT_ID`, `BOKIO_CLIENT_SECRET`

## Company

```
bokio company                          # Get company information
```

## Customers

```
bokio customers list [flags]
bokio customers get <id>
bokio customers create --name <name> [--type company|private] [--email <email>] [--phone <phone>] [--address <addr>] [--city <city>] [--zip-code <zip>] [--country <cc>] [--org-number <num>] [--from-file <path>]
bokio customers update <id> --from-file <path>
bokio customers delete <id>
```

## Suppliers

```
bokio suppliers list [flags]          # filter fields: name, vatNumber, orgNumber
bokio suppliers get <id>
bokio suppliers create --name <name> [--org-number <num>] [--vat-number <num>] [--currency <code>] [--address <addr>] [--city <city>] [--zip-code <zip>] [--country <cc>] [--bankgiro <num> | --plusgiro <num>] [--from-file <path>]
bokio suppliers update <id> [--name <name>] [--org-number <num>] [--vat-number <num>] [--currency <code>] [--bankgiro <num> | --plusgiro <num>] [--from-file <path>]
```

`update` fetches the supplier first and applies changes on top, so unset fields are kept. Requires the `suppliers:write` scope.

## Items

```
bokio items list [flags]
bokio items get <id>
bokio items create --description <desc> [--unit-price <price>] [--unit <unit>] [--vat-rate <rate>] [--from-file <path>]
bokio items update <id> --from-file <path>
bokio items delete <id>
```

## Invoices

```
bokio invoices list [flags]
bokio invoices get <id>
bokio invoices create --customer-id <id> [--invoice-date <YYYY-MM-DD>] [--due-date <YYYY-MM-DD>] [--currency <code>] [--from-file <path>]
bokio invoices update <id> --from-file <path>
bokio invoices publish <id>
bokio invoices record <id> [--account-number <num>] [--payment-date <YYYY-MM-DD>]
```

### Invoice line items

```
bokio invoices line-items add <invoiceId> [--description <desc>] [--quantity <qty>] [--unit-price <price>] [--vat-rate <rate>] [--account-number <num>] [--from-file <path>]
bokio invoices line-items update <invoiceId> <lineItemId> --from-file <path>
bokio invoices line-items delete <invoiceId> <lineItemId>
```

### Invoice attachments

```
bokio invoices attachments list <invoiceId>
bokio invoices attachments add <invoiceId> --file <path>
bokio invoices attachments delete <invoiceId> <attachmentId>
```

### Invoice payments

```
bokio invoices payments list <invoiceId>
bokio invoices payments add <invoiceId> --amount <amount> --payment-date <YYYY-MM-DD> --account-number <num>
```

### Invoice settlements

```
bokio invoices settlements add <invoiceId> --amount <amount> --settle-date <YYYY-MM-DD> --account-number <num>
```

## Journal entries

```
bokio journal-entries list [flags]
bokio journal-entries get <id>
bokio journal-entries create --from-file <path>
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

## Credit notes

```
bokio credit-notes list [flags]
bokio credit-notes get <id>
bokio credit-notes update <id> --from-file <path>
bokio credit-notes publish <id>
bokio credit-notes record <id> [--account-number <num>] [--payment-date <YYYY-MM-DD>]
```

## Uploads

```
bokio uploads list [flags]
bokio uploads get <id>
bokio uploads create --file <path> [--description <desc>] [--journal-entry-id <id>]
bokio uploads download <id> [-O <output-path>]
```

Supported formats: PDF, JPEG, PNG

## Bank payments

```
bokio bank-payments list [flags]
bokio bank-payments get <id>
bokio bank-payments create --from-file <path>
```

## Chart of accounts

```
bokio accounts list
bokio accounts get <accountNumber>
```

## Fiscal years

```
bokio fiscal-years list
bokio fiscal-years get <id>
```

## SIE export

```
bokio sie download <fiscalYearId> [-O <output-path>]
```

Downloads a SIE file (Swedish standard accounting export format) for the given fiscal year.

## Connections

```
bokio connections list [flags]
bokio connections get <id>
bokio connections delete <id>
```

Note: Connections use the General API and do not require `--company-id`.

## List flags (all list commands)

| Flag | Short | Description |
|------|-------|-------------|
| `--page` | | Page number (default: 1) |
| `--page-size` | | Items per page, max 100 (default: 25) |
| `--query` | `-q` | Filter expression |
| `--all` | | Fetch all pages |

## Filter operators

| Operator | Meaning |
|----------|---------|
| `==` | Equals |
| `!=` | Not equals |
| `>` | Greater than |
| `<` | Less than |
| `>=` | Greater or equal |
| `<=` | Less or equal |
| `~` | Contains (strings, case-sensitive) |
| `&&` | AND |
| `\|\|` | OR |

## API rate limits

200 requests per 60-second rolling window per token. The CLI automatically retries on 429 responses.
