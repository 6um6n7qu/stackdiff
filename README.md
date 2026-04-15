# stackdiff

> CLI tool to compare two running service configs and highlight environment drift

---

## Installation

```bash
go install github.com/yourusername/stackdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/stackdiff.git
cd stackdiff
go build -o stackdiff .
```

---

## Usage

Compare two running service configurations and surface environment drift:

```bash
stackdiff --source production --target staging
```

Compare specific services across environments:

```bash
stackdiff --source prod:api-service --target staging:api-service --output json
```

**Example output:**

```
[DRIFT DETECTED]
  KEY                  PRODUCTION        STAGING
  LOG_LEVEL            info              debug
  MAX_CONNECTIONS      100               50
  CACHE_TTL            3600              missing
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--source` | Source environment or service | required |
| `--target` | Target environment or service | required |
| `--output` | Output format: `text`, `json`, `yaml` | `text` |
| `--ignore` | Comma-separated keys to ignore | none |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)