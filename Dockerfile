# syntax=docker/dockerfile:1

FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/go-vault ./cmd/go-vault

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

USER govault
WORKDIR /data

ENV GOVAULT_STORAGE_TYPE=LOCAL \
    GOVAULT_STORAGE_DEST=/data/backups \
    GOVAULT_RUNTIME_META_DB_PATH=/data/meta/go-vault.db \
    GOVAULT_RUNTIME_TEMP_DIR=/tmp/go-vault \
    GOVAULT_API_ADDR=:8080 \
    GOVAULT_SCHEDULE_CRON="0 2 * * *" \
    GOVAULT_SCHEDULE_TIMEZONE=UTC \
    GOVAULT_RETENTION_DAILY=7 \
    GOVAULT_RETENTION_WEEKLY=4 \
    GOVAULT_RETENTION_MONTHLY=12 \
    GOVAULT_DB_SSLMODE=disable \
    TZ=UTC

VOLUME ["/data/backups", "/data/meta"]
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1

CMD ["go-vault", "serve"]
