# Contributing

Thanks for helping improve go-vault.

## How to contribute

1. Open a [GitHub issue](https://github.com/sboy99/go-vault/issues) for bugs or feature ideas before large changes.
2. Fork the repository and create a branch from `main`.
3. Make your change with focused commits.
4. Open a pull request that describes the problem and the fix.

## Project layout

| Path | What it is |
|---|---|
| Repo root | Go module (`go-vault` CLI and `go-vault-server`) |
| `ui/` | Next.js dashboard |

## Local development

Start the demo stack:

```bash
docker compose up -d
```

- Dashboard: http://127.0.0.1:3000
- API: http://127.0.0.1:8080

For Go work on the host, use `make build` and the binaries under `bin/`. For UI work alone, run `npm install` and `npm run dev` inside `ui/`.

## Guidelines

- Keep pull requests small and reviewable.
- Match existing patterns in the area you touch.
- Do not commit secrets, tokens, or real database passwords.
- Prefer tests for bug fixes and new behavior when practical.
