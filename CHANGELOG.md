# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0]

### Added

- PostgreSQL backup and restore via official `pg_dump` / `pg_restore` binaries
- HTTP API and in-process cron scheduler (`go-vault-server`)
- One-shot CLI for setup, backup create/list, and restore (`go-vault`)
- Next.js dashboard for overview, backups, jobs, schedule, and settings
- Multi-arch Docker image (`sboy99/go-vault`) shipping CLI, API, and dashboard
- Local disk and S3 storage backends with GFS retention
