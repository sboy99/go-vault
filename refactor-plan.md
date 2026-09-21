# go-vault refactor plan

Cleanup plan for the package layout: fixing repo hygiene issues, resolving
naming collisions, and tightening the `internal/` vs `pkg/` split. Ordered
from lowest risk to highest, with optional steps clearly marked.

Run `go build ./...`, `go vet ./...`, and `go test ./...` after **every**
step, not just at the end — each step should be its own commit so a broken
build is easy to bisect or revert.

## 1. Stop tracking build/runtime artifacts

`bin/go-vault` (a compiled binary) and `go-vault.db` (a live BoltDB file)
are currently committed. Neither should be version-controlled.

```bash
cat >> .gitignore <<'EOF'
/bin/
*.db
EOF

git rm --cached bin/go-vault go-vault.db
git commit -m "chore: stop tracking build output and runtime db"
```

Zero code risk — do this first.

## 2. Move `pkg/utils` into `internal/utils`

`pkg/` implies "safe for other projects to import." A generic-named utils
grab-bag (`crypto.go`, `time.go`, `reflect.go`, etc.) is app-specific and
belongs under `internal/`.

```bash
git mv pkg/utils internal/utils
grep -rl 'go-vault/pkg/utils' --include='*.go' . | xargs sed -i 's#go-vault/pkg/utils#go-vault/internal/utils#g'
go build ./...
```

Pure path rename, no logic touched — safe to do early and keeps this noisy
import-path diff separate from the riskier merge in step 3.

## 3. Merge the two `config` packages

There's a top-level `config/` (`config.go`, `validator.go`) and a separate
`internal/config/service.go`. Two packages named `config` forces readers to
figure out which one an import refers to every time.

```bash
git mv config/config.go internal/config/config.go
git mv config/validator.go internal/config/validator.go
git mv config/config_test.go internal/config/config_test.go
rmdir config
```

Resolve any name clashes between the moved files and the existing
`service.go`, update imports repo-wide, then:

```bash
go build ./...
go test ./...
```

This is the highest-risk step since it's an actual package merge, not just
a path rename — review the diff carefully before committing.

## 4. Rename `internal/cmd` to `internal/cli`

`cmd/go-vault/main.go` (the entrypoint) and `internal/cmd/` (CLI command
wiring) both have "cmd" in the name for unrelated things.

```bash
git mv internal/cmd internal/cli
grep -rl 'go-vault/internal/cmd' --include='*.go' . | xargs sed -i 's#go-vault/internal/cmd#go-vault/internal/cli#g'
go build ./...
```

Small and mechanical. Pairs naturally with the existing `internal/ui`.

## 5. Rename `internal/database` to `internal/db` (optional)

```bash
git mv internal/database internal/db
grep -rl 'go-vault/internal/database' --include='*.go' . | xargs sed -i 's#go-vault/internal/database#go-vault/internal/db#g'
go build ./...
```

Only worth doing if you want the "are we connected to Postgres"
(`internal/db`) vs "how we dump/restore Postgres" (`internal/engine`) split
to read more clearly at a glance. Skip if you'd rather stop after step 4.

## 6. Split `internal/utils`'s contents by purpose (optional, later)

Once utils is safely under `internal/` (step 2), break it into smaller,
purpose-named packages instead of one grab-bag file set:

- `crypto.go` → `internal/xcrypto`
- `time.go` → `internal/xtime`
- `string_case.go`, `parse.go`, `struct.go`, `reflect.go`, `json.go`,
  `file_system.go` → group or split as makes sense

Do this incrementally, one file/package at a time, each as its own commit.
Not urgent — only worth doing when you're next touching that code anyway,
not as a dedicated pass.

## 7. Verify and tag

Once everything's green:

```bash
go build ./...
go vet ./...
go test ./...
git tag pre-refactor <commit-before-step-1>   # clean rollback point
```

## End state

```
go-vault/
├── cmd/go-vault/main.go
├── internal/
│   ├── alert/
│   ├── api/
│   ├── backup/
│   ├── cli/            # renamed from internal/cmd
│   ├── config/         # merged: struct + validator + service
│   ├── db/              # renamed from database (optional)
│   ├── engine/
│   ├── job/
│   ├── meta/
│   ├── metrics/
│   ├── retention/
│   ├── scheduler/
│   ├── storage/
│   ├── ui/
│   └── utils/           # moved from pkg/utils
├── pkg/
│   ├── boltdb/          # genuinely generic, fine to keep public
│   └── logger/          # genuinely generic, fine to keep public
├── compose.yml
├── Dockerfile
├── Makefile
├── go.mod / go.sum
└── README.md
```

`pkg/boltdb` and `pkg/logger` are the only two packages that earn their
place in `pkg/` — generic enough that another project could import them
standalone. Everything else is app-specific and belongs under `internal/`.