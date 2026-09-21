# go-vault refactor plan

**Status: completed** on branch `refactor/v1` (tag `pre-refactor` is the rollback point).

Cleanup plan for the package layout: fixing repo hygiene issues, resolving
naming collisions, and tightening the `internal/` vs `pkg/` split. Ordered
from lowest risk to highest.

Run `go build ./...`, `go vet ./...`, and `go test ./...` after **every**
step, not just at the end — each step should be its own commit so a broken
build is easy to bisect or revert.

## Completed steps

1. ~~Stop tracking build/runtime artifacts~~ — `.gitignore` already covered `bin/` and `*.db`
2. ~~Move `pkg/utils` into `internal/utils`~~
3. ~~Merge root `config/` into `internal/config`~~ — setup wizard lives in `internal/setup` to avoid a config↔ui import cycle
4. ~~Rename `internal/cmd` to `internal/cli`~~
5. ~~Delete `internal/database`~~ (was unused; not renamed to `internal/db`)
6. ~~Thin `internal/utils`~~ — deleted unused files; kept crypto/time/struct/string_case
7. ~~Extract backup pipeline helpers~~ (`startMeta`, `dumpAndSave`, `verifyDump`, `failBackup`, `BackupAndPrune`)
8. ~~Extract engine `waitOrKill` / binary picker~~
9. ~~Extract setup and prompt helpers~~
10. ~~Cobra `RunE` + serve lifecycle split~~
11. ~~Fix `decodeJobList` + unit tests~~
12. ~~Flatten storage factory~~ (`Backend` interface)

## End state

```
go-vault/
├── cmd/go-vault/main.go
├── internal/
│   ├── alert/
│   ├── api/
│   ├── backup/
│   ├── cli/            # renamed from internal/cmd
│   ├── config/         # merged: struct + validator
│   ├── engine/
│   ├── job/
│   ├── meta/
│   ├── metrics/
│   ├── retention/
│   ├── scheduler/
│   ├── setup/          # interactive wizard (avoids config↔ui cycle)
│   ├── storage/
│   ├── ui/
│   └── utils/          # moved from pkg/utils (thinned)
├── pkg/
│   ├── boltdb/
│   └── logger/
├── compose.yml
├── Dockerfile
├── Makefile
├── go.mod / go.sum
└── README.md
```

`pkg/boltdb` and `pkg/logger` remain the only public packages.
Everything else is app-specific under `internal/`.
