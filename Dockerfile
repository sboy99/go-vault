# syntax=docker/dockerfile:1

FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/go-vault ./cmd/cli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/go-vault-server ./cmd/server

FROM debian:bookworm-slim
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

RUN useradd -r -u 10001 -m -d /home/govault govault \
    && mkdir -p /data/backups /data/meta /tmp/go-vault \
    && chown -R govault:govault /data /tmp/go-vault

COPY --from=build /out/go-vault /usr/local/bin/go-vault
COPY --from=build /out/go-vault-server /usr/local/bin/go-vault-server

USER govault
WORKDIR /data

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

VOLUME ["/data/backups", "/data/meta"]
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1

CMD ["go-vault-server"]
