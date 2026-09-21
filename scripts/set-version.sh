#!/usr/bin/env bash
# Stamp the release version into package.json, the CLI default, and UI mock fixtures.
# Usage: scripts/set-version.sh v1.2.3   or   scripts/set-version.sh 1.2.3
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <version>" >&2
  exit 1
fi

raw="$1"
# Strip one leading v (v1.2.3 → 1.2.3).
ver="${raw#v}"

if [[ ! "$ver" =~ ^[0-9]+\.[0-9]+\.[0-9]+([.-].*)?$ ]]; then
  echo "invalid semver: $raw" >&2
  exit 1
fi

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

# UI package + lockfile (no git tag; allow no-op when already at this version).
(
  cd ui
  npm version "$ver" --no-git-tag-version --allow-same-version
)

# CLI / app.version default source.
version_go="$root/internal/version/version.go"
if [[ ! -f "$version_go" ]]; then
  echo "missing $version_go" >&2
  exit 1
fi
# Portable in-place edit (BSD and GNU sed).
tmp="$(mktemp)"
sed -E "s/var Version = \"[^\"]+\"/var Version = \"$ver\"/" "$version_go" >"$tmp"
mv "$tmp" "$version_go"

# Demo config fixture shown when the API is unreachable.
fixtures="$root/ui/lib/mock/fixtures.ts"
if [[ ! -f "$fixtures" ]]; then
  echo "missing $fixtures" >&2
  exit 1
fi
tmp="$(mktemp)"
sed -E "s/(app: \{ name: \"go-vault\", version: )\"[^\"]+\"/\1\"$ver\"/" "$fixtures" >"$tmp"
mv "$tmp" "$fixtures"

echo "set version to $ver"
