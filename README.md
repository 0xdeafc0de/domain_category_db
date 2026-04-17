# Domain Category DB

A lightweight and extensible domain categorization lookup DB in Go. It supports map, radix-tree, and hash-based indexing for fast domain lookups.

The example application in `main.go` fetches public blocklists from remote URLs, categorizes domains, and serves lookup requests via a REST API.

## Features
- Loads domain blocklists from remote URLs or cached local files.
- Categorizes domains into configurable categories (e.g., abuse, phishing).
- Fast in-memory lookup using map, radix tree, or hashed storage.
- Caches recent lookups with an LRU cache.
- Exposes a RESTful API to check domain categories.
- Easy configuration via `config.json`.

## Project Structure
```bash
.
├── config
│   └── config.go
├── config.json
├── db
│   ├── db.go
│   ├── db_test.go
│   └── db_benchmark_test.go
├── domains
│   └── domains.go
├── main.go
├── Makefile
└── rest
    └── rest.go
```

## Configuration
All categories and data sources are configured in `config.json`:

```json
{
  "dbstore_path": "/path/to/local/cache",
  "categories": [
    {
      "url": "https://blocklistproject.github.io/Lists/alt-version/abuse-nl.txt",
      "category": "abuse"
    },
    {
      "url": "https://blocklistproject.github.io/Lists/alt-version/phishing-nl.txt",
      "category": "phishing"
    }
  ]
}
```

- Domains are cached locally after first download.
- Categories are assigned a unique internal ID starting from 101.

## Build and Run

```bash
make run
```

You can also build directly with:

```bash
go build -o domaindb main.go
./domaindb
```

## Testing

Run the full test suite with:

```bash
go test ./...
```

The `db` package includes unit tests for lookup normalization and file loading, plus benchmark helpers in `db/db_benchmark_test.go`. Some fixture-driven integration tests skip automatically when the local blocklist files are unavailable.

## REST API
The REST API is exposed on localhost:8081 by default.
```bash
GET /info
GET /lookup?domain=<domain>
```

### Example Request using curl
```bash
curl http://127.0.0.1:8081/info
curl "http://127.0.0.1:8081/lookup?domain=facebook.com"
```

