# Go-Vault

Production PostgreSQL backup service. Uses the official `pg_dump` / `pg_restore` binaries, runs on a cron schedule inside Docker, enforces GFS retention (7 daily / 4 weekly / 12 monthly), and exposes an HTTP API to list, download, trigger, and restore backups.

## Features

- Custom-format dumps (`pg_dump -Fc`) with `pg_restore --list` verification
- Streaming upload to local disk or S3 (aws-sdk-go-v2)
- GFS retention with timezone-aware day/week/month boundaries
- Scheduler + single-flight job runner (no overlapping dumps)
- HTTP API on a private Docker network (optional bearer token)
- Prometheus metrics including `govault_backup_last_success_timestamp`

## Quick start (Docker Compose)

```bash
docker compose up -d --build
docker compose exec go-vault curl -s http://127.0.0.1:8080/healthz
docker compose exec go-vault curl -s -X POST http://127.0.0.1:8080/v1/backups
docker compose exec go-vault curl -s http://127.0.0.1:8080/v1/backups
```

The API port is **not** published to the host by default. Attach clients to the `backup_net` network.

## Configuration

Env-first (`GO_VAULT_` prefix). Optional `config.yml` for local CLI use.

| Variable | Default | Description |
|---|---|---|
| `GO_VAULT_DB_HOST` | | Postgres host |
| `GO_VAULT_DB_PORT` | `5432` | Postgres port |
| `GO_VAULT_DB_NAME` | | Database name |
| `GO_VAULT_DB_USERNAME` | | Username |
| `GO_VAULT_DB_PASSWORD` | | Password (or `GO_VAULT_DB_PASSWORD_FILE`) |
| `GO_VAULT_DB_SSLMODE` | `require` | libpq sslmode |
| `GO_VAULT_STORAGE_TYPE` | `LOCAL` | `LOCAL` or `CLOUD` |
| `GO_VAULT_STORAGE_DEST` | `./backups` | Local backup directory |
| `GO_VAULT_SCHEDULE_CRON` | `0 2 * * *` | Daily backup cron |
| `GO_VAULT_SCHEDULE_TIMEZONE` | `UTC` | Cron timezone |
| `GO_VAULT_RETENTION_DAILY` | `7` | Keep newest backup per day for N days |
| `GO_VAULT_RETENTION_WEEKLY` | `4` | Keep newest backup per ISO week for N weeks |
| `GO_VAULT_RETENTION_MONTHLY` | `12` | Keep newest backup per month for N months |
| `GO_VAULT_API_ADDR` | `:8080` | HTTP listen address |
| `GO_VAULT_API_TOKEN` | | Optional bearer token |
| `GO_VAULT_RUNTIME_META_DB_PATH` | `./go-vault.db` | BoltDB path |
| `GO_VAULT_RUNTIME_TEMP_DIR` | `/tmp/go-vault` | Spool for verify/restore |
| `GO_VAULT_RUNTIME_ALERT_WEBHOOK` | | Optional failure webhook URL |

## CLI

```bash
make build
./bin/go-vault setup
./bin/go-vault backup create
./bin/go-vault backup list
./bin/go-vault backup restore <backup_id_or_name>
./bin/go-vault serve
```

## HTTP API

Envelope: `{ "success": bool, "data": ..., "error": ... }`

| Method | Path | Description |
|---|---|---|
| GET | `/healthz` | Liveness |
| GET | `/readyz` | Readiness (DB ping) |
| GET | `/metrics` | Prometheus |
| GET | `/v1/backups` | List backups |
| GET | `/v1/backups/{id}` | Backup metadata |
| GET | `/v1/backups/{id}/download` | Download artifact |
| POST | `/v1/backups` | Trigger backup (+ prune) |
| POST | `/v1/restores` | Restore (`backup_id` + `confirm` = db name) |
| GET | `/v1/jobs/{id}` | Job status |

## Requirements

- Go 1.25+
- `pg_dump` / `pg_restore` with client major >= server major (image ships 15/16/17)

## Releasing

Images are published to Docker Hub as [`sboy99/go-vault`](https://hub.docker.com/r/sboy99/go-vault) when a semver git tag is pushed.

```bash
git tag v0.1.0
git push origin v0.1.0
# → docker pull sboy99/go-vault:0.1.0
#    also tags: 0.1, 0, latest
```

Required GitHub repository secrets (Settings → Secrets and variables → Actions):

| Secret | Description |
|---|---|
| `DOCKERHUB_TOKEN` | Docker Hub access token for user `sboy99` |

## Safety notes

- Restore uses `pg_restore --clean --if-exists` and can leave a partial database on failure. Take a fresh backup before restoring when possible.
- Keep the API on a private network. Set `GO_VAULT_API_TOKEN` if the port is reachable more broadly.
- Credentials live in env / secrets files, not in the dump artifact metadata beyond what Postgres requires at dump time.
