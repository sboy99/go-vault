#!/bin/sh
set -eu

# One-shot CLI / custom command: skip the UI supervisor.
if [ "$#" -gt 0 ] && [ "$1" != "go-vault-server" ]; then
  exec "$@"
fi

# Default UI → API wiring inside the same container.
: "${GOVAULT_API_URL:=http://127.0.0.1:8080}"
: "${GO_VAULT_UI_PORT:=3000}"
export GOVAULT_API_URL
export HOSTNAME=0.0.0.0
export PORT="$GO_VAULT_UI_PORT"

# Reuse the Go API bearer token when the UI token is unset.
if [ -z "${GOVAULT_API_TOKEN:-}" ] && [ -n "${GO_VAULT_API_TOKEN:-}" ]; then
  export GOVAULT_API_TOKEN="$GO_VAULT_API_TOKEN"
fi

UI_PID=""
SERVER_PID=""

term_handler() {
  if [ -n "$SERVER_PID" ]; then
    kill -TERM "$SERVER_PID" 2>/dev/null || true
  fi
  if [ -n "$UI_PID" ]; then
    kill -TERM "$UI_PID" 2>/dev/null || true
  fi
  wait "$SERVER_PID" 2>/dev/null || true
  wait "$UI_PID" 2>/dev/null || true
  exit 0
}

trap term_handler TERM INT

cd /opt/govault-ui
node server.js &
UI_PID=$!

go-vault-server &
SERVER_PID=$!

# Exit if either child dies (propagate non-zero status when possible).
while true; do
  if ! kill -0 "$UI_PID" 2>/dev/null; then
    wait "$UI_PID" || status=$?
    kill -TERM "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
    exit "${status:-1}"
  fi
  if ! kill -0 "$SERVER_PID" 2>/dev/null; then
    wait "$SERVER_PID" || status=$?
    kill -TERM "$UI_PID" 2>/dev/null || true
    wait "$UI_PID" 2>/dev/null || true
    exit "${status:-1}"
  fi
  sleep 1
done
