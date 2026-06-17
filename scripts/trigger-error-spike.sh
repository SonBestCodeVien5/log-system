#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="${LOG_FILE:-./logs/demo-node/app.log}"
SERVICE_NAME="${SERVICE_NAME:-demo-node}"
COUNT="${1:-20}"
MESSAGE="${MESSAGE:-INCIDENT_REPLAY: controlled error spike for alert testing}"

if ! [[ "$COUNT" =~ ^[1-9][0-9]*$ ]]; then
  echo "usage: $0 [positive-count]" >&2
  exit 1
fi

if ! [[ "$SERVICE_NAME" =~ ^[A-Za-z0-9._-]+$ ]]; then
  echo "SERVICE_NAME must contain only letters, numbers, dot, underscore, or hyphen" >&2
  exit 1
fi

if [[ "$MESSAGE" == *'"'* || "$MESSAGE" == *'\'* || "$MESSAGE" == *$'\n'* ]]; then
  echo "MESSAGE must not contain quotes, backslashes, or newlines" >&2
  exit 1
fi

BATCH_ID="$(date -u +%Y%m%dT%H%M%SZ)-$$"
LINES_FILE="$(mktemp)"
trap 'rm -f "$LINES_FILE"' EXIT

for i in $(seq 1 "$COUNT"); do
  timestamp="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
  printf '{"timestamp":"%s","level":"ERROR","service":"%s","message":"%s","metadata":{"source":"incident-replay","batch_id":"%s","sequence":%d,"count":%d}}\n' \
    "$timestamp" "$SERVICE_NAME" "$MESSAGE" "$BATCH_ID" "$i" "$COUNT" >> "$LINES_FILE"
done

log_dir="$(dirname "$LOG_FILE")"
if mkdir -p "$log_dir" && { cat "$LINES_FILE" >> "$LOG_FILE"; } 2>/dev/null; then
  echo "wrote $COUNT ERROR logs to $LOG_FILE for service=$SERVICE_NAME batch_id=$BATCH_ID"
  exit 0
fi

container_service=""
case "$LOG_FILE" in
  ./logs/demo-node/app.log|logs/demo-node/app.log)
    container_service="demo-node"
    ;;
  ./logs/demo-go/app.log|logs/demo-go/app.log)
    container_service="demo-go"
    ;;
esac

if [[ -n "$container_service" ]] && command -v docker >/dev/null 2>&1; then
  if docker compose ps -q "$container_service" >/dev/null 2>&1 &&
    docker compose exec -T "$container_service" sh -c 'cat >> /var/log/app/app.log' < "$LINES_FILE"; then
    echo "wrote $COUNT ERROR logs to $LOG_FILE via container=$container_service service=$SERVICE_NAME batch_id=$BATCH_ID"
    exit 0
  fi
fi

echo "failed to write $COUNT ERROR logs to $LOG_FILE" >&2
echo "hint: start docker compose or choose a writable LOG_FILE path" >&2
exit 1
