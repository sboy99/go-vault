# syntax=docker/dockerfile:1
#
# go-vault — PostgreSQL backup image
# ----------------------------------
# Scheduled pg_dump / pg_restore, GFS retention, an HTTP API, and a web dashboard.
#
#   docker pull sboy99/go-vault:latest
#   API        http://<host>:8080     (go-vault-server, cron included)
#   Dashboard  http://<host>:3000
#
# One-shot CLI (does not start the scheduler or the dashboard):
#   docker run --rm sboy99/go-vault:latest go-vault backup list
#
# Stage 1 (build):   compile the CLI and server binaries
# Stage 2 (ui):      production Next.js standalone bundle
# Stage 3 (runtime): slim Debian with pg clients 15/16/17, API, and dashboard
#
#   go-vault         CLI (setup, backup create|list|restore)
#   go-vault-server  HTTP API + cron scheduler
#   /opt/govault-ui  Next.js dashboard (GO_VAULT_UI_PORT, default 3000)

# =============================================================================
# Stage 1 — Go build
# =============================================================================

FROM golang:1.25-bookworm AS build
WORKDIR /src

# Release tag stamped into binaries (overridden by docker build --build-arg VERSION=).
ARG VERSION=0.1.0
# Buildx sets these per target platform (linux/amd64, linux/arm64, …).
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Step 1.1: Copy module files first so dependency download is cached
#           independently of source changes.
COPY go.mod go.sum ./

# Step 1.2: Download Go module dependencies into the module cache.
RUN go mod download

# Step 1.3: Copy the full source tree into the build context.
COPY . .

# Step 1.4: Compile both binaries as static Linux artifacts for the target arch.
#           CGO_ENABLED=0       — no libc dependency (runs on slim Debian)
#           -trimpath           — reproducible builds (no host paths in binary)
#           -ldflags="-s -w"    — strip symbol/DWARF tables (smaller binary)
#           -X …Version         — stamp release version into the CLI / app default
#           Outputs land in /out for the runtime stage to copy.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w -X github.com/sboy99/go-vault/internal/version.Version=${VERSION}" \
      -o /out/go-vault ./cmd/cli && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w -X github.com/sboy99/go-vault/internal/version.Version=${VERSION}" \
      -o /out/go-vault-server ./cmd/server

# =============================================================================
# Stage 2 — UI build (discarded after copy; only standalone output is kept)
# =============================================================================

FROM node:22-bookworm-slim AS ui
WORKDIR /ui

ENV NEXT_TELEMETRY_DISABLED=1

# Node 22 bundles npm 10. This lockfile is written by npm 11, and npm 10's
# `npm ci` rejects it. This image's npm 10.9.8 can upgrade itself.
RUN npm install -g npm@11.6.2

# Step 2.1: Install deps from the lockfile (include devDeps — next/tsc/tailwind compile the app).
COPY ui/package.json ui/package-lock.json ./
RUN npm ci

# Step 2.2: Copy UI source and produce the standalone production bundle.
COPY ui/ ./
RUN NODE_ENV=production npm run build

# =============================================================================
# Stage 3 — Runtime
# =============================================================================

FROM debian:bookworm-slim

LABEL org.opencontainers.image.title="go-vault" \
      org.opencontainers.image.description="PostgreSQL backups with pg_dump, GFS retention, an HTTP API, and a web dashboard." \
      org.opencontainers.image.source="https://github.com/sboy99/go-vault" \
      org.opencontainers.image.url="https://github.com/sboy99/go-vault" \
      org.opencontainers.image.documentation="https://github.com/sboy99/go-vault#readme"

# Step 3.1: Install OS packages needed at runtime, then clean apt lists.
#           - ca-certificates / curl / gnupg / lsb-release — TLS + PGDG repo setup
#           - libstdc++6 — required by the Node binary copied from the ui stage
#           - postgresql-client-15/16/17 — pg_dump / pg_restore matching server majors
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates curl gnupg lsb-release libstdc++6 \
    && curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc \
      | gpg --dearmor -o /usr/share/keyrings/postgresql.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/postgresql.gpg] http://apt.postgresql.org/pub/repos/apt bookworm-pgdg main" \
      > /etc/apt/sources.list.d/pgdg.list \
    && apt-get update && apt-get install -y --no-install-recommends \
      postgresql-client-15 \
      postgresql-client-16 \
      postgresql-client-17 \
    && rm -rf /var/lib/apt/lists/*

# Step 3.2: Create a non-root user and data directories for backups + BoltDB meta.
RUN useradd -r -u 10001 -m -d /home/govault govault \
    && mkdir -p /data/backups /data/meta /tmp/go-vault /opt/govault-ui \
    && chown -R govault:govault /data /tmp/go-vault

# Step 3.3: Copy Go binaries and the Node runtime binary (not the Node image).
COPY --from=build /out/go-vault /usr/local/bin/go-vault
COPY --from=build /out/go-vault-server /usr/local/bin/go-vault-server
COPY --from=ui /usr/local/bin/node /usr/local/bin/node

# Step 3.4: Copy only the traced Next.js standalone production output.
#           Built with WORKDIR /ui, so server.js sits at standalone/ root.
COPY --from=ui --chown=govault:govault /ui/.next/standalone/ /opt/govault-ui/
COPY --from=ui --chown=govault:govault /ui/.next/static /opt/govault-ui/.next/static

# Step 3.5: Supervisor that runs API + UI together (CLI overrides via CMD).
COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Step 3.6: Drop privileges and set the working directory to the data volume root.
USER govault
WORKDIR /data

# Step 3.7: Default runtime configuration (override via env / compose / k8s).
#           Storage is local under /data; scheduler runs daily at 02:00 UTC;
#           GFS retention keeps 7 daily / 4 weekly / 12 monthly backups.
ENV GO_VAULT_STORAGE_TYPE=LOCAL \
    GO_VAULT_STORAGE_DEST=/data/backups \
    GO_VAULT_RUNTIME_META_DB_PATH=/data/meta/go-vault.db \
    GO_VAULT_RUNTIME_TEMP_DIR=/tmp/go-vault \
    GO_VAULT_API_ADDR=:8080 \
    GO_VAULT_UI_PORT=3000 \
    GOVAULT_API_URL=http://127.0.0.1:8080 \
    GO_VAULT_SCHEDULE_CRON="0 2 * * *" \
    GO_VAULT_SCHEDULE_TIMEZONE=UTC \
    GO_VAULT_RETENTION_DAILY=7 \
    GO_VAULT_RETENTION_WEEKLY=4 \
    GO_VAULT_RETENTION_MONTHLY=12 \
    GO_VAULT_DB_SSLMODE=disable \
    TZ=UTC

# Step 3.8: Declare persistent volumes and listen ports (API + UI).
VOLUME ["/data/backups", "/data/meta"]
EXPOSE 8080 3000

# Step 3.9: Liveness probe against the HTTP health endpoint.
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1

# Step 3.10: Default process runs API + UI via the entrypoint.
#            Override with `go-vault` for one-shot CLI use inside the container.
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["go-vault-server"]
