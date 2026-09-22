#!/usr/bin/env bash
# Insert a Keep a Changelog section for a release tag.
#
# Existing version headings are left unchanged, so hand-written notes stay.
# When ## [Unreleased] contains bullets, that section is promoted to the release.
# Otherwise the section is built from conventional commits since the previous
# stable semver tag (see _classify_commit).
#
# Usage: scripts/update-changelog.sh v1.2.3
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

# _require_tag exits 1 when the git tag does not exist.
# $1 = tag name (e.g. v1.2.3)
_require_tag() {
  local tag="$1"

  if ! git rev-parse --verify --quiet "refs/tags/${tag}^{commit}" >/dev/null; then
    echo "tag ${tag} not found" >&2
    exit 1
  fi
}

# _escape_regex escapes characters that are special inside a basic/extended regex.
# $1 = literal string
# Prints the escaped string.
_escape_regex() {
  printf '%s' "$1" | sed -E 's/[][(){}.^$*+?|\\]/\\&/g'
}

# _changelog_has_version returns 0 when CHANGELOG.md already has ## [version].
# $1 = bare version
# $2 = path to CHANGELOG.md
_changelog_has_version() {
  local ver="$1"
  local changelog="$2"
  local escaped

  escaped="$(_escape_regex "$ver")"
  grep -Eq "^## \\[${escaped}]([[:space:]]|$)" "$changelog"
}

# _tag_date returns the commit date of a tag as YYYY-MM-DD.
# $1 = tag name
_tag_date() {
  local tag="$1"
  local date

  date="$(git log -1 --format=%cs "$tag")"
  if [[ -z "$date" ]]; then
    date="$(git log -1 --format=%ci "$tag" | cut -c1-10)"
  fi
  printf '%s\n' "$date"
}

# _previous_tag finds the previous semver tag already merged into $tag.
# For a stable release (X.Y.Z), pre-release tags are skipped.
# $1 = current tag (e.g. v1.2.3)
# $2 = bare version
# Prints the previous tag, or an empty line when none exists.
_previous_tag() {
  local tag="$1"
  local ver="$2"
  local stable=0
  local candidate

  if [[ "$ver" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    stable=1
  fi

  while IFS= read -r candidate; do
    if [[ "$candidate" == "$tag" ]]; then
      continue
    fi
    if [[ "$stable" -eq 1 && ! "$candidate" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      continue
    fi
    printf '%s\n' "$candidate"
    return
  done < <(git tag --list 'v[0-9]*.[0-9]*.[0-9]*' --merged "$tag" --sort=-v:refname)

  printf '\n'
}

# _github_repo_url resolves https://github.com/owner/repo from the environment
# or the origin remote. Prints the URL, or an empty line when unknown.
_github_repo_url() {
  if [[ -n "${GITHUB_REPOSITORY:-}" ]]; then
    printf 'https://github.com/%s\n' "$GITHUB_REPOSITORY"
    return
  fi

  local origin
  origin="$(git remote get-url origin 2>/dev/null || true)"
  origin="${origin%.git}"

  if [[ "$origin" =~ ^git@github.com:(.+)$ ]]; then
    printf 'https://github.com/%s\n' "${BASH_REMATCH[1]}"
    return
  fi
  if [[ "$origin" =~ ^ssh://git@github.com/(.+)$ ]]; then
    printf 'https://github.com/%s\n' "${BASH_REMATCH[1]}"
    return
  fi
  if [[ "$origin" =~ ^https://github.com/(.+)$ ]]; then
    printf 'https://github.com/%s\n' "${BASH_REMATCH[1]}"
    return
  fi

  printf '\n'
}

# _append_compare_link appends a Keep a Changelog footer link for this version.
# $1 = bare version
# $2 = current tag
# $3 = previous tag (may be empty)
# $4 = path to CHANGELOG.md
_append_compare_link() {
  local ver="$1"
  local tag="$2"
  local previous="$3"
  local changelog="$4"
  local base url line escaped

  base="$(_github_repo_url)"
  if [[ -z "$base" ]]; then
    echo "could not resolve GitHub repo; skipping compare link" >&2
    return
  fi

  if [[ -n "$previous" ]]; then
    url="${base}/compare/${previous}...${tag}"
  else
    url="${base}/releases/tag/${tag}"
  fi

  line="[${ver}]: ${url}"
  escaped="$(_escape_regex "$ver")"
  if grep -Eq "^\\[${escaped}]:[[:space:]]" "$changelog"; then
    return
  fi

  if [[ -s "$changelog" && -n "$(tail -c 1 "$changelog")" ]]; then
    printf '\n' >>"$changelog"
  fi
  if ! grep -Eq '^\[[^]]+]:[[:space:]]' "$changelog"; then
    printf '\n' >>"$changelog"
  fi
  printf '%s\n' "$line" >>"$changelog"
}

# _unreleased_has_notes returns 0 when ## [Unreleased] contains bullet lines.
# $1 = path to CHANGELOG.md
_unreleased_has_notes() {
  local changelog="$1"

  awk '
    /^## \[Unreleased\][[:space:]]*$/ { in_u=1; next }
    /^## / && in_u { exit 1 }
    in_u && /^- / { found=1; exit 0 }
    END { exit found ? 0 : 1 }
  ' "$changelog"
}

# _promote_unreleased renames ## [Unreleased] to ## [version] - date and
# inserts a fresh empty Unreleased heading above it.
# $1 = bare version
# $2 = release date (YYYY-MM-DD)
# $3 = path to CHANGELOG.md
_promote_unreleased() {
  local ver="$1"
  local date="$2"
  local changelog="$3"
  local heading out

  heading="## [${ver}] - ${date}"
  out="$(mktemp)"
  awk -v heading="$heading" '
    /^## \[Unreleased\][[:space:]]*$/ && !done {
      print "## [Unreleased]"
      print ""
      print heading
      done=1
      next
    }
    { print }
  ' "$changelog" >"$out"
  mv "$out" "$changelog"
  echo "promoted [Unreleased] to ${ver}"
}

# _commit_is_omitted returns 0 for conventional types that stay out of the notes.
# $1 = conventional commit type (lowercase)
_commit_is_omitted() {
  case "$1" in
    chore|docs|doc|test|tests|ci|style|build|refactor|release) return 0 ;;
    *) return 1 ;;
  esac
}

# _body_has_breaking_change returns 0 when the commit body declares a breaking change.
# $1 = commit body
_body_has_breaking_change() {
  printf '%s' "$1" | grep -qi 'BREAKING[- ]CHANGE:'
}

# _capitalize uppercases the first character of a string.
# $1 = text
# Prints the capitalized text.
_capitalize() {
  local msg="$1"
  local first rest

  first="$(printf '%s' "${msg:0:1}" | tr '[:lower:]' '[:upper:]')"
  rest="${msg:1}"
  printf '%s%s' "$first" "$rest"
}

# _classify_commit maps a conventional commit to a Keep a Changelog group.
#
# Mapping:
#   feat     -> added
#   fix      -> fixed
#   perf     -> changed
#   security -> security
#   type! or BREAKING CHANGE footer -> changed (breaking)
#   chore/docs/test/ci/style/build/refactor -> omitted (empty group)
#
# $1 = subject line
# $2 = body
# Prints: "<group>\t<bullet>"  or nothing when the commit is omitted.
_classify_commit() {
  local subject="$1"
  local body="$2"
  # Variable form keeps parentheses inside the regex on bash 3.2.
  local conventional_re='^([A-Za-z]+)(\([^)]+\))?(!)?:[[:space:]]+(.+)$'
  local type="" scope="" bang="" msg group bullet

  msg="$subject"
  if [[ "$subject" =~ $conventional_re ]]; then
    type="$(printf '%s' "${BASH_REMATCH[1]:-}" | tr '[:upper:]' '[:lower:]')"
    scope="${BASH_REMATCH[2]:-}"
    bang="${BASH_REMATCH[3]:-}"
    msg="${BASH_REMATCH[4]:-}"
  fi

  msg="$(printf '%s' "$msg" | sed -E 's/[[:space:]]+/ /g; s/^ //; s/ $//')"
  if [[ -z "$msg" ]]; then
    return
  fi
  msg="$(_capitalize "$msg")"

  if [[ -n "$scope" ]]; then
    scope="${scope#(}"
    scope="${scope%)}"
    bullet="${msg} (${scope})"
  else
    bullet="$msg"
  fi

  group="changed"
  case "$type" in
    feat|feature) group="added" ;;
    fix|bugfix) group="fixed" ;;
    security) group="security" ;;
    perf|revert) group="changed" ;;
  esac

  if [[ -n "$bang" ]] || _body_has_breaking_change "$body"; then
    group="changed"
    bullet="**Breaking:** ${bullet}"
  elif [[ -n "$type" ]] && _commit_is_omitted "$type"; then
    return
  fi

  printf '%s\t%s\n' "$group" "$bullet"
}

# _collect_notes builds bullet files under $workdir from git history.
# $1 = git range (e.g. v0.1.0..v0.2.0)
# $2 = workdir with added/changed/fixed/security files
_collect_notes() {
  local range="$1"
  local workdir="$2"
  local subject body classified group bullet

  while true; do
    IFS= read -r -d '' subject || break
    IFS= read -r -d '' body || true
    # git log writes a newline after each NUL record; it lands on the next subject.
    subject="${subject#$'\n'}"
    if [[ -z "$subject" ]]; then
      continue
    fi

    classified="$(_classify_commit "$subject" "$body")"
    if [[ -z "$classified" ]]; then
      continue
    fi

    group="${classified%%$'\t'*}"
    bullet="${classified#*$'\t'}"

    if grep -Fxq -- "- ${bullet}" "$workdir/$group"; then
      continue
    fi
    printf -- '- %s\n' "$bullet" >>"$workdir/$group"
  done < <(git log --reverse --no-merges --format='%s%x00%b%x00' "$range")
}

# _render_section writes a Keep a Changelog section to $outdir/section.
# $1 = bare version
# $2 = release date
# $3 = workdir with group bullet files
# $4 = output path for the rendered section
_render_section() {
  local ver="$1"
  local date="$2"
  local workdir="$3"
  local out="$4"
  local wrote=0
  local title file

  {
    printf '## [%s] - %s\n\n' "$ver" "$date"
    for title in Added Changed Fixed Security; do
      case "$title" in
        Added) file="$workdir/added" ;;
        Changed) file="$workdir/changed" ;;
        Fixed) file="$workdir/fixed" ;;
        Security) file="$workdir/security" ;;
      esac
      if [[ -s "$file" ]]; then
        printf '### %s\n\n' "$title"
        cat "$file"
        printf '\n'
        wrote=1
      fi
    done
    if [[ "$wrote" -eq 0 ]]; then
      printf '### Changed\n\n- Maintenance release with no user-facing changes.\n\n'
    fi
  } >"$out"
}

# _insert_section places a rendered section after ## [Unreleased], creating
# that heading when the file does not have one yet.
# $1 = path to the rendered section file
# $2 = path to CHANGELOG.md
_insert_section() {
  local section_file="$1"
  local changelog="$2"
  local out

  out="$(mktemp)"
  awk '
    NR == FNR { section = section $0 ORS; next }
    /^## \[Unreleased\][[:space:]]*$/ { in_u = 1 }
    in_u && /^## / && !/^## \[Unreleased\]/ && !inserted {
      printf "%s", section
      inserted = 1
    }
    !in_u && /^## / && !inserted {
      print "## [Unreleased]"
      print ""
      printf "%s", section
      inserted = 1
    }
    { print }
    END {
      if (!inserted) {
        print "## [Unreleased]"
        print ""
        printf "%s", section
      }
    }
  ' "$section_file" "$changelog" >"$out"
  mv "$out" "$changelog"
}

# _generate_section builds notes from conventional commits and inserts them.
# $1 = bare version
# $2 = current tag
# $3 = previous tag (may be empty)
# $4 = release date
# $5 = path to CHANGELOG.md
_generate_section() {
  local ver="$1"
  local tag="$2"
  local previous="$3"
  local date="$4"
  local changelog="$5"
  local workdir range

  workdir="$(mktemp -d)"
  : >"$workdir/added"
  : >"$workdir/changed"
  : >"$workdir/fixed"
  : >"$workdir/security"

  if [[ -n "$previous" ]]; then
    range="${previous}..${tag}"
  else
    range="$tag"
  fi

  _collect_notes "$range" "$workdir"
  _render_section "$ver" "$date" "$workdir" "$workdir/section"
  _insert_section "$workdir/section" "$changelog"
  rm -rf "$workdir"

  if [[ -n "$previous" ]]; then
    echo "wrote changelog ${ver} from ${previous}..${tag}"
  else
    echo "wrote changelog ${ver} from ${tag}"
  fi
}

# main updates CHANGELOG.md for a release tag.
# $@ = CLI args (exactly one version or tag string)
main() {
  if [[ $# -ne 1 ]]; then
    echo "usage: $0 <version>" >&2
    exit 1
  fi

  local ver tag root changelog date previous

  ver="$(_normalize_version "$1")"
  tag="v${ver}"
  root="$(_repo_root)"
  cd "$root"
  changelog="$root/CHANGELOG.md"

  if [[ ! -f "$changelog" ]]; then
    echo "missing $changelog" >&2
    exit 1
  fi

  _require_tag "$tag"

  if _changelog_has_version "$ver" "$changelog"; then
    echo "changelog already has ${ver}"
    exit 0
  fi

  date="$(_tag_date "$tag")"
  previous="$(_previous_tag "$tag" "$ver")"

  if _unreleased_has_notes "$changelog"; then
    _promote_unreleased "$ver" "$date" "$changelog"
  else
    _generate_section "$ver" "$tag" "$previous" "$date" "$changelog"
  fi

  _append_compare_link "$ver" "$tag" "$previous" "$changelog"
}

main "$@"
