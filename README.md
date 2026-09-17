# Go-Vault

CLI for backing up and restoring a database. Run `go-vault setup` once, then create, list, and restore dumps from the terminal.

Today only **PostgreSQL** is implemented. MySQL and MongoDB show up in setup, but they are not wired yet.

## What it does

1. Walks you through database and storage settings and writes `config.yml`.
2. Connects to Postgres, dumps schema plus data as SQL (no `pg_dump` binary required).
3. Saves the dump locally (`./backups` by default) or to an S3 bucket.
4. Records backup metadata in a local BoltDB file (`go-vault.db`).
5. Restores a named dump back into the configured database.

Dump files are named `<unix-timestamp>_<db-name>_backup.sql`.

## Requirements

- Go 1.22+
- A reachable PostgreSQL database
- For cloud storage: an S3 bucket (or any S3-compatible endpoint)

## Install

```bash
git clone https://github.com/sboy99/go-vault.git
cd go-vault
make build
```

The binary lands at `./bin/go-vault`.

```bash
make install   # same as build
make test
make clean
```

Or without Make:

```bash
go build -o ./bin/go-vault ./cmd/go-vault
```

### Docker

```bash
docker build -t go-vault .
docker run --rm -it go-vault
```

The image runs `./go-vault` with no args, which prints help. Mount a working directory if you want `config.yml` and backups to persist.

## Usage

All commands except `setup` need a `config.yml` in the current directory.

```bash
# Interactive config (DB + storage)
go-vault setup

# Create a dump and store it
go-vault backup create

# Show recent backups (up to 15)
go-vault backup list

# Restore a dump by filename
go-vault backup restore 1710000000_postgres_backup.sql
```

`backup restore` has shell completion for dump names from metadata. Generate completions with `make completions`.

## Configuration

`go-vault setup` writes `config.yml`. You can also create it by hand:

```yaml
app:
  name: go-vault
  version: 0.0.1
db:
  name: postgres
  type: POSTGRESQL
  host: localhost
  port: 5432
  username: postgres
  password: secret
storage:
  type: LOCAL          # LOCAL or CLOUD
  dest: ./backups
  cloud:
    type: AWS          # AWS only for now
    aws:
      region: ap-south-1
      bucket_name: my-backups
      access_key_id: ...
      access_key_secret: ...
      endpoint: default  # or a custom S3-compatible URL
```

Set `storage.cloud.aws.endpoint` to something other than `default` for MinIO or other S3-compatible stores (path-style addressing is enabled).

`config.yml` and `go-vault.db` are gitignored. They contain credentials and backup history — keep them out of source control.

Create the local destination before the first backup if it does not already exist:

```bash
mkdir -p backups
```

## How a backup is built

Postgres dumps are generated in-process against `information_schema` and table data:

- schemas (non-system)
- extensions
- `CREATE TABLE` for public tables
- sequences
- primary keys
- `COPY ... FROM stdin` row data

Restore reads that SQL file and executes it against the configured database.

## Project layout

```
cmd/go-vault/     CLI entrypoint
internal/cmd/     Cobra commands
internal/config/  Interactive setup
internal/database Postgres adapter
internal/storage  Local disk + AWS S3
internal/meta     Backup metadata (BoltDB)
internal/ui       Prompts and table output
pkg/pg_dump/      SQL dump generator
config/           Config load/save/validate
```

## Current limits

- PostgreSQL only. MySQL and MongoDB are placeholders.
- GCP cloud storage is a placeholder. AWS S3 works; custom endpoints are supported.
- Postgres connections use `sslmode=disable`.
- Credentials are stored in plaintext in `config.yml`.
- Backups are not encrypted at rest by go-vault.
- Cloud delete is not implemented.
- The dump covers public tables, schemas, extensions, sequences, and primary keys. Indexes, foreign keys, views, functions, and roles are not dumped.
