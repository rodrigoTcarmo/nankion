# Nankion

A Go CLI tool that parses OFX (Open Financial Exchange) bank statements and syncs them to Notion databases. Perfect for personal finance tracking and automating bank statement management.

## Features

- **OFX Parsing**: Reads standard OFX files from banks (supports both bank accounts and credit cards statements)
- **Notion Integration**: Automatically uploads transactions as pages to your Notion database
- **Schema Management**: Auto-creates required database properties if they don't exist
- **Batch Upload**: Process multiple OFX files from a folder in one command
- **Dual Logging**: Outputs logs to both console and timestamped log files
- **Full OFX Support**: Handles all standard OFX transaction types (CREDIT, DEBIT, PAYMENT, XFER, etc.)

## Installation

### Prerequisites

- Go 1.22.1 or higher
- A Notion integration with access to your target database

### Build from Source

```bash
git clone https://github.com/rodrigoTcarmo/nankion.git
cd nankion
go install ./cmd/nankion
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `NOTION_SECRET` | Your Notion integration secret/API key | Yes |

To get your Notion secret:
1. Go to [Notion Integrations](https://www.notion.so/my-integrations)
2. Create a new integration
3. Copy the "Internal Integration Secret"
4. Share your target database with the integration

```bash
export NOTION_SECRET="secret_your_notion_integration_key"
```

## Usage

### Upload a Single Statement

```bash
nankion report upload <database-id> <path-to-ofx-file>
```

**Example:**
```bash
nankion report upload abc123def456 ./statements/september-2025.ofx
```

### Upload Multiple Statements

```bash
nankion report upload-statements <database-id> <folder-path>
```

**Example:**
```bash
nankion report upload-statements abc123def456 ./statements/
```

This will process all `.ofx` files in the specified folder.

### Database Commands

**Get database info:**
```bash
nankion database get <database-id>
```

**List database properties:**
```bash
nankion database get-properties <database-id>
```

## Notion Database Schema

Nankion automatically creates the following properties in your Notion database if they don't exist:

| Property | Type | Description |
|----------|------|-------------|
| `Name` | Title | Auto-generated: `{Operation}-{Destination}-{TransactionID}` |
| `Transaction Date` | Date | Date the transaction was posted |
| `Operation` | Select | Transaction type (CREDIT, DEBIT, PAYMENT, etc.) |
| `Destination` | Rich Text | Payee/merchant name |
| `Amount` | Number | Transaction amount (Brazilian Real format) |
| `Memo` | Rich Text | Additional transaction details (e.g., PIX info) |
| `Transaction ID` | Rich Text | Unique transaction identifier (FITID) |

### Supported Operation Types

- `CREDIT` - Generic credit (money in)
- `DEBIT` - Generic debit (money out)
- `PAYMENT` - Electronic payment
- `XFER` - Transfer
- `DEP` - Deposit
- `ATM` - ATM transaction
- `POS` - Point of sale
- `CHECK` - Check
- `FEE` - Bank fee
- `INT` - Interest earned
- `DIV` - Dividend
- `DIRECTDEP` - Direct deposit
- `DIRECTDEBIT` - Merchant initiated debit
- `REPEATPMT` - Repeating payment
- `SRVCHG` - Service charge
- `CASH` - Cash withdrawal
- `OTHER` - Other

## OFX Format Support

Nankion parses the standard OFX/QFX format used by most banks:

```
OFXHEADER:100
<OFX>
  <BANKMSGSRSV1>
    <STMTTRNRS>
      <STMTRS>
        <BANKTRANLIST>
          <STMTTRN>
            <TRNTYPE>CREDIT
            <DTPOSTED>20250915
            <TRNAMT>1500.00
            <FITID>unique-transaction-id
            <NAME>Employer Inc
            <MEMO>Salary deposit
          </STMTTRN>
        </BANKTRANLIST>
        <LEDGERBAL>
          <BALAMT>5000.00
          <DTASOF>20250930
        </LEDGERBAL>
      </STMTRS>
    </STMTTRNRS>
  </BANKMSGSRSV1>
</OFX>
```

Supports both:
- Bank statements (`BANKMSGSRSV1`)
- Credit card statements (`CREDITCARDMSGSRSV1`)

## Project Structure

```
nankion/
├── cmd/
│   ├── nankion/          # Main CLI application
│   │   └── main.go
├── pkg/
│   ├── log/              # Logging utilities (dual output)
│   ├── models/
│   │   ├── ofx/          # OFX transaction type definitions
│   │   └── properties/   # Notion property builders
│   ├── notion/
│   │   ├── client/       # Notion API client wrapper
│   │   ├── database/     # Database operations
│   │   └── page/         # Page creation and management
│   └── statement/        # OFX parsing and statement handling
├── go.mod
├── go.sum
└── README.md
```

## Dependencies

- [ofxgo](https://github.com/aclindsa/ofxgo) - OFX file parsing
- [go-notion](https://github.com/dstotijn/go-notion) - Notion API client
- [cobra](https://github.com/spf13/cobra) - CLI framework

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build main CLI
go build -o nankion ./cmd/nankion

```

### Vendoring

Dependencies are vendored for reproducible builds:

```bash
go mod vendor
```

## Logging

Nankion creates daily log files in the format `nankion-DD-MM-YYYY.log` in the current working directory. Logs are written to both stdout and the log file simultaneously.

## Finding Your Database ID

1. Open your Notion database in a browser
2. The URL will look like: `https://www.notion.so/workspace/abc123def456?v=...`
3. The database ID is: `abc123def456` (32 characters after your workspace name)

Alternatively, you can get the database ID from:
- The "Copy link" option in Notion (extract from URL)
- The Notion API when querying your workspace

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Nankion** - Because managing finances should be as simple as running a command.
