# go-vault UI

SolidJS + Vite dashboard for [go-vault](../README.md), styled after the pg-view design system.

## Setup

```bash
cp .env.example .env
npm install
npm run dev
```

The Vite `/api` proxy forwards to `VITE_API_PROXY` (default `http://localhost:8080`).

## Scripts

| Command | Description |
|---|---|
| `npm run dev` | Start Vite dev server |
| `npm run build` | Typecheck + production build |
| `npm test` | Run Vitest suite |
| `npm run test:coverage` | Vitest with coverage |

## Stack

- SolidJS 1.9 + `@solidjs/router`
- Tailwind CSS v4 (pg-view / Nuxt UI v4 token layer)
- Kobalte headless primitives
- Global store with repository-backed slices
