# Go-Vault

PostgreSQL backup service. It dumps with the official `pg_dump` / `pg_restore` binaries, stores artifacts on local disk or S3, keeps a cycle rollup retention set (7 daily / 4 weekly / 12 monthly by default), and runs the schedule inside one process.

The Docker image ships three things:

| Piece | What it is | Default |
|---|---|---|
| `go-vault` | CLI for setup, create, list, and restore | one-shot; does not start the scheduler |
| `go-vault-server` | HTTP API + cron scheduler | listens on `:8080` |
| Dashboard | Next.js UI in the same container | [http://127.0.0.1:3000](http://127.0.0.1:3000) |

Images: [`sboy99/go-vault`](https://hub.docker.com/r/sboy99/go-vault) (`linux/amd64` and `linux/arm64`).

## Quick start

`compose.yml` starts Postgres 16 and go-vault, and publishes the API and the dashboard.

```bash
docker compose up -d
```

- Dashboard: [http://127.0.0.1:3000](http://127.0.0.1:3000)
- API: [http://127.0.0.1:8080](http://127.0.0.1:8080)

```bash
curl -s http://127.0.0.1:8080/healthz
curl -s -X POST http://127.0.0.1:8080/v1/backups
curl -s http://127.0.0.1:8080/v1/jobs
curl -s http://127.0.0.1:8080/v1/backups
```

`POST /v1/backups` returns a job. Poll `GET /v1/jobs/{id}` until `status` is `succeeded` or `failed`. Only one backup or restore runs at a time.

Backups land in the `go_vault_backups` volume (`/data/backups`). Backup and job metadata live in BoltDB on `go_vault_meta` (`/data/meta/go-vault.db`).

The compose file is a local demo: database password `postgres`, TLS off, API open on the host with no token. For anything else, set a real password, `GO_VAULT_DB_SSLMODE=require`, and `GO_VAULT_API_TOKEN`.

### Your own Postgres

Skip the bundled database and point the image at an existing server. The container must be able to reach `GO_VAULT_DB_HOST`.

```bash
docker run -d --name go-vault \
  -p 8080:8080 -p 3000:3000 \
  -v go_vault_backups:/data/backups \
  -v go_vault_meta:/data/meta \
  -e GO_VAULT_DB_HOST=db.internal \
  -e GO_VAULT_DB_PORT=5432 \
  -e GO_VAULT_DB_NAME=app \
  -e GO_VAULT_DB_USERNAME=backup \
  -e GO_VAULT_DB_PASSWORD=secret \
  -e GO_VAULT_DB_SSLMODE=require \
  -e GO_VAULT_API_TOKEN=change-me \
  sboy99/go-vault:latest
```

Open [http://127.0.0.1:3000](http://127.0.0.1:3000). With `GO_VAULT_API_TOKEN` set, the dashboard reuses that token for its own API calls. Callers outside the container send `Authorization: Bearer change-me`. `/healthz` and `/readyz` stay open so probes work without the token.

## Dashboard

The default container process starts the API and the UI together.

| Page | What you see |
|---|---|
| Overview | Backup counts, size, recent backups, recent jobs |
| Jobs | Async backup and restore jobs |
| Schedule | Cron expression, next run, rollup keep/prune preview |
| Settings | Redacted database, storage, and retention config |

The UI talks to the API at `GOVAULT_API_URL` (default `http://127.0.0.1:8080` inside the image). It does not show passwords or cloud secrets.

If the API is down, pages fall back to demo fixtures and show a demo-data banner.

## CLI

Two binaries share the same config and services. `go-vault` does not run the scheduler. `go-vault-server` does.

### On the host

Requires Go 1.25+ and `pg_dump` / `pg_restore` / `psql` matching the server major when possible. A newer client can dump an older server, but restoring into that server needs a matching (or carefully filtered) client.

```bash
make build
./bin/go-vault setup
./bin/go-vault backup create
./bin/go-vault backup list
./bin/go-vault backup restore <backup_id_or_name>
./bin/go-vault-server
```

`setup` is interactive. It asks for the database and, if you pick cloud storage, the S3 bucket, then writes `config.yml` in the working directory. Other commands read that file and overlay `GO_VAULT_*` environment variables.

`backup restore` asks you to type the database name before it runs `pg_restore`.

### Inside the container

Any command other than `go-vault-server` skips the UI and runs that command.

```bash
docker compose run --rm go-vault go-vault backup list
docker compose run --rm go-vault go-vault backup create
docker compose run --rm go-vault go-vault backup restore <backup_id_or_name>
```

`backup create` in the CLI runs the dump and retention prune in the foreground and prints the backup id. The HTTP `POST /v1/backups` path does the same work as a background job.

## HTTP API

Responses use `{ "success": true, "data": ..., "error": null }`. Failures set `success` to false and put the message in `error`.

When `GO_VAULT_API_TOKEN` is set, send `Authorization: Bearer <token>` on every route except `GET /healthz` and `GET /readyz`.

| Method | Path | Description |
|---|---|---|
| GET | `/healthz` | Liveness |
| GET | `/readyz` | Readiness (database ping) |
| GET | `/metrics` | Prometheus |
| GET | `/v1/backups` | List backups (`limit`, `offset`, `status`) |
| GET | `/v1/backups/{id}` | One backup. `id` is the backup id or the artifact name |
| GET | `/v1/backups/{id}/download` | Download the custom-format dump |
| POST | `/v1/backups` | Start a backup job, then prune |
| POST | `/v1/restores` | Start a restore job |
| GET | `/v1/jobs` | List jobs (`limit`, `offset`) |
| GET | `/v1/jobs/{id}` | One job |
| GET | `/v1/config` | Redacted config, including `next_run` |
| GET | `/v1/stats` | Counts, total size, and whether a job is running |
| GET | `/v1/retention` | Backups the rollup policy would keep and prune |

Trigger a backup:

```bash
curl -s -X POST http://127.0.0.1:8080/v1/backups \
  -H "Authorization: Bearer $GO_VAULT_API_TOKEN"
```

Restore. `confirm` must be the configured database name. `backup_id` may be the id or the artifact name.

```bash
curl -s -X POST http://127.0.0.1:8080/v1/restores \
  -H "Authorization: Bearer $GO_VAULT_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"backup_id":"<id-or-name>","confirm":"app"}'
```

Download:

```bash
curl -fsS -OJ http://127.0.0.1:8080/v1/backups/<id>/download \
  -H "Authorization: Bearer $GO_VAULT_API_TOKEN"
```

Prometheus metric `govault_backup_last_success_timestamp` is the unix time of the last successful backup.

## Configuration

Environment variables use the `GO_VAULT_` prefix and override `config.yml`. Nested keys become underscores: `db.host` is `GO_VAULT_DB_HOST`.

Secrets can be passed as files. If the `*_FILE` variable is set, its contents replace the matching value:

| File variable | Replaces |
|---|---|
| `GO_VAULT_DB_PASSWORD_FILE` | `GO_VAULT_DB_PASSWORD` |
| `GO_VAULT_API_TOKEN_FILE` | `GO_VAULT_API_TOKEN` |
| `GO_VAULT_STORAGE_CLOUD_AWS_ACCESS_KEY_SECRET_FILE` | `GO_VAULT_STORAGE_CLOUD_AWS_ACCESS_KEY_SECRET` |

### Database

| Variable | Default | Description |
|---|---|---|
| `GO_VAULT_DB_HOST` | | Postgres host |
| `GO_VAULT_DB_PORT` | `5432` | Postgres port |
| `GO_VAULT_DB_NAME` | | Database name |
| `GO_VAULT_DB_USERNAME` | | Username |
| `GO_VAULT_DB_PASSWORD` | | Password |
| `GO_VAULT_DB_TYPE` | `POSTGRESQL` | Only `POSTGRESQL` is implemented |
| `GO_VAULT_DB_SSLMODE` | `require` | libpq sslmode. The image sets `disable` so the bundled compose Postgres works |

`GO_VAULT_DB_HOST`, name, username, and password are required. The server exits on startup if they are missing.

### Storage

| Variable | Default | Description |
|---|---|---|
| `GO_VAULT_STORAGE_TYPE` | `LOCAL` | `LOCAL` or `CLOUD` |
| `GO_VAULT_STORAGE_DEST` | `./backups` | Local directory. The image uses `/data/backups` |

Cloud storage is S3 (or an S3-compatible endpoint):

| Variable | Description |
|---|---|
| `GO_VAULT_STORAGE_CLOUD_TYPE` | `AWS` |
| `GO_VAULT_STORAGE_CLOUD_AWS_REGION` | Region |
| `GO_VAULT_STORAGE_CLOUD_AWS_BUCKET_NAME` | Bucket |
| `GO_VAULT_STORAGE_CLOUD_AWS_ACCESS_KEY_ID` | Access key |
| `GO_VAULT_STORAGE_CLOUD_AWS_ACCESS_KEY_SECRET` | Secret key |
| `GO_VAULT_STORAGE_CLOUD_AWS_ENDPOINT` | Optional custom endpoint. `default` uses AWS |

```bash
GO_VAULT_STORAGE_TYPE=CLOUD \
GO_VAULT_STORAGE_CLOUD_TYPE=AWS \
GO_VAULT_STORAGE_CLOUD_AWS_REGION=us-east-1 \
GO_VAULT_STORAGE_CLOUD_AWS_BUCKET_NAME=my-backups \
GO_VAULT_STORAGE_CLOUD_AWS_ACCESS_KEY_ID=AKIA... \
GO_VAULT_STORAGE_CLOUD_AWS_ACCESS_KEY_SECRET=... \
./bin/go-vault-server
```

### Schedule and retention

| Variable | Default | Description |
|---|---|---|
| `GO_VAULT_SCHEDULE_CRON` | `0 2 * * *` | Five-field cron. Default is 02:00 every day |
| `GO_VAULT_SCHEDULE_TIMEZONE` | `UTC` | Timezone for cron and for week/month cycle boundaries |
| `GO_VAULT_RETENTION_DAILY` | `7` | Cap on daily backups kept in the open week |
| `GO_VAULT_RETENTION_WEEKLY` | `4` | Cap on weekly backups (Saturday of each closed week) kept in the open month |
| `GO_VAULT_RETENTION_MONTHLY` | `12` | Cap on monthly backups (last calendar day of each closed month) |

Each value must be at least 1. Weeks run Sunday–Saturday. When a week closes, its Saturday backup is promoted to weekly and the rest of that week is pruned (newest in the week if Saturday is missing). When a month closes, its last-day backup is promoted to monthly and the rest of that month is pruned. The newest successful backup is always kept. Failed backups are not pruned. Prune runs after each successful backup.

The scheduler starts with `go-vault-server`. It will not start a new dump while a backup or restore job is already running.

### API, UI, and runtime

| Variable | Default | Description |
|---|---|---|
| `GO_VAULT_API_ADDR` | `:8080` | HTTP listen address |
| `GO_VAULT_API_TOKEN` | empty | Bearer token. Empty disables auth |
| `GO_VAULT_UI_PORT` | `3000` | Dashboard port inside the image |
| `GOVAULT_API_URL` | `http://127.0.0.1:8080` | Where the dashboard sends API calls |
| `GOVAULT_API_TOKEN` | falls back to `GO_VAULT_API_TOKEN` | Token the dashboard uses. Server-side only |
| `GO_VAULT_RUNTIME_META_DB_PATH` | `./go-vault.db` | BoltDB path. The image uses `/data/meta/go-vault.db` |
| `GO_VAULT_RUNTIME_TEMP_DIR` | `/tmp/go-vault` | Spool for verify and restore |
| `GO_VAULT_RUNTIME_DUMP_TIMEOUT` | `2h` | `pg_dump` timeout |
| `GO_VAULT_RUNTIME_RESTORE_TIMEOUT` | `2h` | `pg_restore` timeout |
| `GO_VAULT_RUNTIME_SHUTDOWN_TIMEOUT` | `5m` | Grace period on SIGINT / SIGTERM |
| `GO_VAULT_RUNTIME_RESTORE_JOBS` | `4` | Parallel jobs passed to `pg_restore` |
| `GO_VAULT_RUNTIME_ALERT_WEBHOOK` | empty | POST a JSON body on backup or restore failure |

Webhook body:

```json
{"event":"backup_failed","message":"...","time":"2026-09-22T00:00:00Z"}
```

`event` is `backup_failed` or `restore_failed`.

## What a backup does

1. `pg_dump -Fc` writes a custom-format dump.
2. `pg_restore --list` checks that the archive can be read.
3. The file is streamed to local disk or S3. Metadata (id, name, size, sha256, Postgres version, status) is stored in BoltDB.
4. Retention deletes artifacts that the cycle rollup no longer keeps.

On startup the server reconciles BoltDB with the files it can still see in storage.

Restore runs `pg_restore --clean --if-exists` into the configured database. A failed restore can leave that database partially rewritten. Take a fresh backup before restoring when you can.

## Build from source

```bash
make build          # ./bin/go-vault and ./bin/go-vault-server
make test
docker build -t go-vault:dev .
```

The image installs PostgreSQL client 15, 16, and 17 and prefers the client that matches the server major. Archives written by a newer `pg_dump` are restored through a filtered SQL path (single-threaded via `psql`) so older servers do not see unsupported settings such as `transaction_timeout`.

UI-only development (the dashboard against a running API, or demo fixtures) is documented in [`ui/README.md`](ui/README.md).

## Releasing

A semver git tag publishes `sboy99/go-vault` to Docker Hub.

```bash
git tag v0.1.0
git push origin v0.1.0
# docker pull sboy99/go-vault:0.1.0
# also tagged 0.1, 0, and latest
```

GitHub Actions secret `DOCKERHUB_TOKEN` is a Docker Hub access token for user `sboy99`.

## Safety

- Keep the API and the dashboard on a private network. Set `GO_VAULT_API_TOKEN` if port 8080 or 3000 is reachable beyond that network.
- The dashboard proxies backup, restore, and download. Do not put it on the public internet without an auth gate in front of it.
- Passwords and cloud secrets belong in the environment or `*_FILE` secrets, not in the image.
