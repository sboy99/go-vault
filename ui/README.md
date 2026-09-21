# go-vault UI

Next.js dashboard for the go-vault PostgreSQL backup API.

## Setup

```bash
cp .env.example .env.local
npm install
npm run dev
```

Environment:

| Variable | Description |
| --- | --- |
| `GOVAULT_API_URL` | Base URL of the go-vault HTTP API (default `http://127.0.0.1:8080`) |
| `GOVAULT_API_TOKEN` | Bearer token matching `api.token` in the Go server config (server-side only) |
| `GOVAULT_METRICS_PUBLIC_URL` | Optional browser-reachable Prometheus URL |
| `GOVAULT_USE_MOCKS` | Set to `1` to force demo fixtures |

When the API is unreachable, the UI falls back to typed mock fixtures and shows a "demo data" banner.

## Security

- `GOVAULT_API_TOKEN` is read only in server components / route handlers (`import "server-only"`).
- Do not expose this UI on the public internet without a reverse-proxy auth gate. The BFF proxies privileged API operations (backup, restore, download).
- Prefer binding `next start` / the reverse proxy to localhost or a private network.
- Mutating `/api/*` requests are blocked when `Sec-Fetch-Site` indicates a cross-site caller.
