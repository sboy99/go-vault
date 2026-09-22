#!/usr/bin/env bash
# Stamp the release version into the UI package, the CLI default, and UI mock fixtures.
#
# Usage: scripts/set-version.sh v1.2.3
#    or: scripts/set-version.sh 1.2.3
#
# Entry point: main (bottom of file).
set -euo pipefail

# _normalize_version strips one leading "v" and rejects non-semver input.
# $1 = raw version (e.g. v1.2.3 or 1.2.3)
# Prints the bare version (e.g. 1.2.3) or exits 1.
_normalize_version() {
  local raw="$1"
  local ver="${raw#v}"

  if [[ ! "$ver" =~ ^[0-9]+\.[0-9]+\.[0-9]+([.-].*)?$ ]]; then
    echo "invalid semver: $raw" >&2
    exit 1
  fi

  printf '%s\n' "$ver"
}

# _repo_root resolves the go-vault repository root from this script's location.
# Prints an absolute path.
_repo_root() {
  cd "$(dirname "$0")/.." && pwd
}

# _replace_pattern rewrites a file with sed via a temp file (BSD and GNU safe).
# $1 = absolute file path
# $2 = sed -E expression
_replace_pattern() {
  local file="$1"
  local expression="$2"
  local tmp

  tmp="$(mktemp)"
  sed -E "$expression" "$file" >"$tmp"
  mv "$tmp" "$file"
}

# _stamp_ui_package sets ui/package.json (+ lockfile) to the given version.
# Does not create a git tag. No-ops when already at that version.
# $1 = bare version
# $2 = repository root
_stamp_ui_package() {
  local ver="$1"
  local root="$2"

  (
    cd "$root/ui"
    npm version "$ver" --no-git-tag-version --allow-same-version
  )
}

# _stamp_go_version sets var Version in internal/version/version.go.
# $1 = bare version
# $2 = repository root
_stamp_go_version() {
  local ver="$1"
  local root="$2"
  local version_go="$root/internal/version/version.go"

  if [[ ! -f "$version_go" ]]; then
    echo "missing $version_go" >&2
    exit 1
  fi

  _replace_pattern "$version_go" \
    "s/var Version = \"[^\"]+\"/var Version = \"$ver\"/"
}

# _stamp_mock_fixtures sets app.version in the UI demo fixture.
# $1 = bare version
# $2 = repository root
_stamp_mock_fixtures() {
  local ver="$1"
  local root="$2"
  local fixtures="$root/ui/lib/mock/fixtures.ts"

  if [[ ! -f "$fixtures" ]]; then
    echo "missing $fixtures" >&2
    exit 1
  fi

  _replace_pattern "$fixtures" \
    "s/(app: \{ name: \"go-vault\", version: )\"[^\"]+\"/\1\"$ver\"/"
}

# main stamps every version source for a release.
# $@ = CLI args (exactly one version string)
main() {
  if [[ $# -ne 1 ]]; then
    echo "usage: $0 <version>" >&2
    exit 1
  fi

  local ver root
  ver="$(_normalize_version "$1")"
  root="$(_repo_root)"
  cd "$root"

  _stamp_ui_package "$ver" "$root"
  _stamp_go_version "$ver" "$root"
  _stamp_mock_fixtures "$ver" "$root"

  echo "set version to $ver"
}

main "$@"
