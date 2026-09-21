# syntax=docker/dockerfile:1
#
# go-vault multi-stage image
# --------------------------
# Stage 1 (build):  compile the CLI and server binaries
# Stage 2 (runtime): slim Debian image with pg_dump/pg_restore clients
#
# Binaries produced:
#   go-vault        — CLI (setup, backup create|list|restore)
#   go-vault-server — HTTP API + cron scheduler (default CMD)

# =============================================================================
# Stage 1 — Build
# =============================================================================

FROM golang:1.25-bookworm AS build
WORKDIR /src

# Step 1.1: Copy module files first so dependency download is cached
#           independently of source changes.
COPY go.mod go.sum ./

# Step 1.2: Download Go module dependencies into the module cache.
RUN go mod download

# Step 1.3: Copy the full source tree into the build context.
COPY . .

# Step 1.4: Cross-compile both binaries as static Linux amd64 artifacts.
#           CGO_ENABLED=0       — no libc dependency (runs on slim Debian)
#           -trimpath           — reproducible builds (no host paths in binary)
#           -ldflags="-s -w"    — strip symbol/DWARF tables (smaller binary)
#           Outputs land in /out for the runtime stage to copy.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/go-vault ./cmd/cli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/go-vault-server ./cmd/server

# =============================================================================
# Stage 2 — Runtime
# =============================================================================

FROM debian:bookworm-slim

# Step 2.1: Install OS packages needed at runtime, then clean apt lists.
#           - ca-certificates / curl / gnupg / lsb-release — TLS + PGDG repo setup
#           - postgresql-client-15/16/17 — pg_dump / pg_restore matching server majors
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates curl gnupg lsb-release \
    && curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc \
      | gpg --dearmor -o /usr/share/keyrings/postgresql.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/postgresql.gpg] http://apt.postgresql.org/pub/repos/apt bookworm-pgdg main" \
      > /etc/apt/sources.list.d/pgdg.list \
    && apt-get update && apt-get install -y --no-install-recommends \
      postgresql-client-15 \
      postgresql-client-16 \
      postgresql-client-17 \
    && rm -rf /var/lib/apt/lists/*

# Step 2.2: Create a non-root user and data directories for backups + BoltDB meta.
RUN useradd -r -u 10001 -m -d /home/govault govault \
    && mkdir -p /data/backups /data/meta /tmp/go-vault \
    && chown -R govault:govault /data /tmp/go-vault

# Step 2.3: Copy compiled binaries from the build stage into PATH.
COPY --from=build /out/go-vault /usr/local/bin/go-vault
COPY --from=build /out/go-vault-server /usr/local/bin/go-vault-server

# Step 2.4: Drop privileges and set the working directory to the data volume root.
USER govault
WORKDIR /data

# Step 2.5: Default runtime configuration (override via env / compose / k8s).
#           Storage is local under /data; scheduler runs daily at 02:00 UTC;
#           GFS retention keeps 7 daily / 4 weekly / 12 monthly backups.
ENV GO_VAULT_STORAGE_TYPE=LOCAL \
    GO_VAULT_STORAGE_DEST=/data/backups \
    GO_VAULT_RUNTIME_META_DB_PATH=/data/meta/go-vault.db \
    GO_VAULT_RUNTIME_TEMP_DIR=/tmp/go-vault \
    GO_VAULT_API_ADDR=:8080 \
    GO_VAULT_SCHEDULE_CRON="0 2 * * *" \
    GO_VAULT_SCHEDULE_TIMEZONE=UTC \
    GO_VAULT_RETENTION_DAILY=7 \
    GO_VAULT_RETENTION_WEEKLY=4 \
    GO_VAULT_RETENTION_MONTHLY=12 \
    GO_VAULT_DB_SSLMODE=disable \
    TZ=UTC

# Step 2.6: Declare persistent volumes and the API listen port.
VOLUME ["/data/backups", "/data/meta"]
EXPOSE 8080

# Step 2.7: Liveness probe against the HTTP health endpoint.
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1

# Step 2.8: Default process is the server (API + cron).
#           Override with `go-vault` for one-shot CLI use inside the container.
CMD ["go-vault-server"]
