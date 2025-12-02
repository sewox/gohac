# GoHAC CMS

A high-performance, AI-native, multi-tenant Headless CMS built with Go.

## Features

- **Single Binary**: Distributed as standalone Community Edition or SaaS Enterprise Edition
- **Block Protocol**: Flexible, schema-less content structure using JSON blocks
- **Multi-Tenancy**: Subdomain-based tenant resolution (Enterprise Edition)
- **High Performance**: Built on Fiber v2 framework
- **Database**: SQLite (default) or PostgreSQL support via build tags

## Architecture

```
gohac/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── core/            # Domain entities and repository interfaces
│   │   ├── domain/      # Domain models (Page, Block, etc.)
│   │   └── repository/  # Repository interfaces
│   ├── adapter/         # Implementations (SQL, HTTP handlers)
│   └── middleware/      # Auth, Tenant Resolution, CORS
└── web/
    ├── admin/           # React Admin Panel source
    └── theme/           # Astro Theme source
```

## Build Tags

- `community`: SQLite, Local Auth, Local FS (default)
- `enterprise`: PostgreSQL, Multi-tenancy, S3 support

## Getting Started

```bash
# Install dependencies
go mod download

# Run server (Community Edition)
go run cmd/server/main.go

# Server will start on http://localhost:3000
```

## License

MIT

