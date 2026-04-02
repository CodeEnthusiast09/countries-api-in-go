# Country Currency & Exchange API

A RESTful API built with Go (Gin) that fetches country data and live exchange rates, stores them in MySQL, and exposes endpoints for querying, filtering, sorting, and image generation.

## Table of Contents

- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Local Development](#local-development)
- [Database Migrations](#database-migrations)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [Deployment (Railway)](#deployment-railway)
- [Architecture Decisions](#architecture-decisions)

---

## Tech Stack

- **Go** with [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- **GORM** — ORM for MySQL
- **Atlas** — versioned database migrations
- **MySQL** — relational database
- **Docker** — containerization for deployment

---

## Project Structure

```
.
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Loads environment variables
│   ├── database/
│   │   ├── database.go          # DB connection via GORM
│   │   └── migrate.go           # Atlas programmatic migration runner (local use)
│   ├── handlers/
│   │   ├── country_handler.go   # HTTP handlers
│   │   └── dtos.go              # Request/response types
│   ├── image/
│   │   └── generator.go         # PNG summary image generation
│   ├── lib/
│   │   └── utils.go             # Shared utilities
│   ├── models/
│   │   ├── country.go           # Country GORM model
│   │   └── external.go          # External API response types
│   ├── router/
│   │   └── router.go            # Route definitions
│   └── services/
│       ├── country_service.go   # Business logic
│       ├── external_service.go  # Calls to RestCountries + Exchange Rate APIs
│       └── types.go             # Service interfaces
├── migrations/
│   ├── 20260401221520_create_countries_table.sql
│   └── atlas.sum                # Atlas integrity checksum file — do not edit manually
├── atlas.hcl                    # Atlas configuration (used for generating migrations locally)
├── Dockerfile                   # Multi-stage build for production
├── entrypoint.sh                # Container startup: runs migrations then starts server
├── railway.toml                 # Railway deployment configuration
└── tools.go                     # Blank import for Atlas GORM provider
```

---

## Prerequisites

- Go 1.26+
- MySQL 8+ (or MariaDB locally)
- [Atlas CLI](https://atlasgo.io/docs) installed
- Docker (for building/running the container locally)

Install Atlas CLI:

```bash
curl -sSf https://atlasgo.sh | sh
```

---

## Local Development

1. Clone the repository and install dependencies:

```bash
git clone https://github.com/CodeEnthusiast09/country-currency-api.git
cd country-currency-api
go mod download
```

2. Create a local MySQL/MariaDB database:

```sql
CREATE DATABASE countries_api;
```

3. Copy and fill in your environment variables:

```bash
cp .env.example .env
```

4. Run migrations (see [Database Migrations](#database-migrations) below)

5. Start the server:

```bash
go run ./cmd/
```

The API will be available at `http://localhost:3000`.

---

## Database Migrations

This project uses [Atlas](https://atlasgo.io) for versioned migrations — the same approach used in production. Never use `AutoMigrate` in production; it cannot roll back and gives you no audit trail of schema changes.

### How Atlas tracks what's been applied

Atlas maintains a table called `atlas_schema_revisions` in your database. Every time you run `atlas migrate apply`, it checks this table to see which `.sql` files have already been applied and only runs the new ones. Think of it as a "receipt" of all migrations that have run.

### atlas.sum — the integrity file

The `migrations/atlas.sum` file is a checksum of every file in the `migrations/` directory. Atlas verifies this before running anything. This protects you from situations where someone edits a migration file that was already applied to a database — which would cause a dangerous mismatch between your schema history and reality.

**Rule:** never edit a migration file that has already been applied to any database. If you need to change the schema, create a new migration file instead.

### Apply migrations locally

```bash
atlas migrate apply \
  --env gorm \
  --url "mysql://user:password@localhost:3306/countries_api"
```

### Generate a new migration after changing a model

After modifying a GORM model in `internal/models/`, generate a new migration file:

```bash
atlas migrate diff <migration_name> --env gorm
```

This compares your GORM models against the current database schema and generates a `.sql` file containing only the difference (the "diff").

### After editing a migration file manually

If you ever need to edit a `.sql` migration file that has **not yet been applied** to any database, you must update the checksum file or Atlas will refuse to run:

```bash
atlas migrate hash
```

This rewrites `atlas.sum` to match the current state of the migrations directory. Always commit both the `.sql` file and the updated `atlas.sum` together.

---

## Environment Variables

| Variable | Description | Example |
|---|---|---|
| `DATABASE_URL` | GORM DSN for MySQL | `user:pass@tcp(host:3306)/db?charset=utf8mb4&parseTime=True&loc=Local` |
| `ATLAS_DATABASE_URL` | Atlas DSN for MySQL | `mysql://user:pass@host:3306/db` |
| `DB_NAME` | Database name | `countries_api` |
| `PORT` | Server port | `3000` |
| `GIN_MODE` | Gin mode | `release` (production) or `debug` (local) |
| `COUNTRIES_API_URL` | RestCountries API URL | `https://restcountries.com/v2/all?fields=...` |
| `EXCHANGE_RATE_API_URL` | Exchange rate API URL | `https://open.er-api.com/v6/latest/USD` |

Note: `DATABASE_URL` and `ATLAS_DATABASE_URL` point to the same database. GORM uses its own DSN format; Atlas uses a standard URL format.

---

## API Endpoints

### POST `/countries/refresh`
Fetches all countries and exchange rates from external APIs and upserts them into the database.

```json
{ "message": "Countries data refreshed successfully", "total": 250 }
```

### GET `/countries`
Returns all countries with optional filtering and sorting.

| Query param | Description | Example |
|---|---|---|
| `region` | Filter by region | `Africa`, `Europe` |
| `currency` | Filter by currency code | `NGN`, `USD` |
| `sort` | Sort order | `gdp_desc`, `population_asc`, `name_asc` |

```bash
GET /countries?region=Africa&sort=population_desc
```

### GET `/countries/:name`
Returns a single country by name.

```bash
GET /countries/Nigeria
```

### DELETE `/countries/:name`
Deletes a country from the database.

### GET `/status`
Returns total country count and last refresh timestamp. Used as the health check endpoint.

### GET `/countries/image`
Returns a generated PNG image summarising the top 5 countries by GDP.

---

## Deployment (Railway)

The production deployment uses Docker. The `Dockerfile` is a multi-stage build:

- **Stage 1 (builder):** uses the full Go image to compile the binary
- **Stage 2 (runtime):** starts from a clean Alpine image, installs only what's needed to run (ca-certificates + Atlas CLI), copies the compiled binary

The `entrypoint.sh` script runs on every container start and:
1. Waits for the database to accept connections (retries up to 30 times)
2. Runs `atlas migrate apply` — only applies migrations that haven't been applied yet
3. Starts the server with `exec /app/server`

This means migrations are always applied **before** the app starts, and if a migration fails the container exits immediately with a visible error — the app never starts in a broken state.

### Railway environment variables

Set these on your app service in Railway. Use Railway's variable reference syntax (`${{MySQL.VAR_NAME}}`) to pull values from your MySQL service rather than hardcoding credentials:

```
DATABASE_URL    = root:${{MySQL.MYSQL_ROOT_PASSWORD}}@tcp(${{MySQL.MYSQLHOST}}:3306)/railway?charset=utf8mb4&parseTime=True&loc=Local
ATLAS_DATABASE_URL = ${{MySQL.MYSQL_URL}}
COUNTRIES_API_URL  = https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies
EXCHANGE_RATE_API_URL = https://open.er-api.com/v6/latest/USD
PORT    = 3000
GIN_MODE = release
```

### First deployment checklist

- [ ] MySQL service added in Railway project
- [ ] All environment variables set on the app service
- [ ] `Dockerfile`, `entrypoint.sh`, `railway.toml` committed and pushed
- [ ] After deployment succeeds, seed data by calling `POST /countries/refresh`

---

## Architecture Decisions

**Why versioned migrations instead of AutoMigrate?**
`gorm.AutoMigrate` is convenient locally but dangerous in production. It cannot drop columns, cannot roll back, and leaves no audit trail. Versioned migrations give you full control: you know exactly what ran, when, and in what order.

**Why is migration logic in `entrypoint.sh` and not in Go code?**
Running migrations inside the app startup (`database.New()`) creates a race condition when you have multiple replicas — they all try to migrate simultaneously. Putting migrations in the entrypoint ensures they run once, sequentially, before the app process starts.

**Why two DB URL formats?**
GORM uses the Go `database/sql` driver DSN format (`user:pass@tcp(host:port)/db`). Atlas uses a standard URL format (`mysql://user:pass@host:port/db`). They connect to the same database — the format difference is just driver convention.

**Why pin exact package versions in the Dockerfile?**
`apk add --no-cache ca-certificates` installs whatever is current at build time. `apk add --no-cache ca-certificates=20250911-r0` locks the exact version. This makes the build reproducible — you get the same binary whether you build today or six months from now.
