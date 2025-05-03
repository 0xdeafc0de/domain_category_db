# Domain Category DB
A lightweight and extensible domain categorization lookup DB in Go. It supports optional use of radix trees and hash-based indexing for high-performance lookups.

The example application (main.go), fetches public blocklists from remote URLs, categorizes domains, and serves lookup requests via a REST API. 

## Features
- Loads domain blocklists from remote URLs or cached local files.
- Categorizes domains into configurable categories (e.g., abuse, phishing).
- Fast in-memory lookup using map, radix tree, or hashed storage.
- Caches recent lookups with an LRU cache (1M entries).
- Exposes a RESTful API to check domain categories.
- Easy configuration via `config.json`.

## Project Structure
```bash
.
├── config
│   └── config.go
├── config.json
├── db
│   ├── db_benchmark_test.go
│   └── db.go
├── domaindb
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
domain_category_db > make
Building domaindb...
go build -o domaindb main.go

domain_category_db > 
domain_category_db > ./domaindb
2025/05/02 14:55:09 DBStore Path =  /tmp/spingal/domain-category-db/db/dbstore
Loaded 435156 domains from /tmp/spingal/domain-category-db/db/dbstore/abuse into category abuse
Loaded 26031 domains from /tmp/spingal/domain-category-db/db/dbstore/drugs into category drugs
Loaded 190222 domains from /tmp/spingal/domain-category-db/db/dbstore/phishing into category phishing
Loaded 22459 domains from /tmp/spingal/domain-category-db/db/dbstore/facebook into category facebook
Loaded 500282 domains from /tmp/spingal/domain-category-db/db/dbstore/porn into category porn
Loaded 2624 domains from /tmp/spingal/domain-category-db/db/dbstore/torrent into category torrent
Loaded 15070 domains from /tmp/spingal/domain-category-db/db/dbstore/tracking into category tracking
Loaded 1904 domains from /tmp/spingal/domain-category-db/db/dbstore/ransomware into category ransomware
2025/05/02 14:55:09 Total 8 categories loaded. Total DB Count = 1193748, Size ~6.58 MB
Starting server on port 8081
Starting pprof server on :8082
```

## REST API
The REST API is exposed on localhost:8081 by default.
```bash
GET /info
GET /lookup?domain=<domain>
```

### Example Request using curl
```bash
 > curl http://127.0.0.1:8081/info
=== Domain Category DB Info ===
Total DB entries         : 1181605
Approximate DB size      : 25298846 bytes
Cached entries (LRU)     : 1
Categories (ID → Name)   :
  101 → abuse
  102 → drugs
  103 → phishing
  104 → facebook
  105 → porn
  106 → torrent
  107 → tracking
  108 → ransomware

> curl "http://127.0.0.1:8081/lookup?domain=facebook.com"
Category: facebook
Lookup Time: 2125 ns
```

