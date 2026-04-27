# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files with namespace filtering

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultpull.git
cd vaultpull
go build -o vaultpull .
```

---

## Usage

Set your Vault address and token, then run `vaultpull` with a path and optional namespace filter:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

# Pull all secrets from a path into a .env file
vaultpull --path secret/myapp --output .env

# Filter secrets by namespace prefix
vaultpull --path secret/myapp --namespace production --output .env.production
```

**Example output (`.env`):**

```
DB_HOST=db.example.com
DB_PASSWORD=supersecret
API_KEY=abc123
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to read from | *(required)* |
| `--output` | Output `.env` file path | `.env` |
| `--namespace` | Filter secrets by namespace prefix | *(none)* |
| `--overwrite` | Overwrite existing `.env` file | `false` |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with a valid token and read access

---

## License

[MIT](LICENSE)