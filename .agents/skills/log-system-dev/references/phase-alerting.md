# Phase: Alerting And WebSocket

Use this context for `api-server/alerting` and alert UI integration.

## Behavior

- Alert engine runs in its own goroutine.
- It periodically counts `ERROR` logs in the sliding window.
- It broadcasts an alert only when count exceeds threshold and cooldown allows it.
- It supports dynamic threshold/config updates without restarting.

## Shared State

Protect mutable shared state with `sync.RWMutex`, including:

- threshold and timing config
- connected WebSocket clients
- deduplication timestamps

## Config

Read these from environment variables with safe defaults:

- `ALERT_THRESHOLD`
- `ALERT_WINDOW_SECONDS`
- `ALERT_COOLDOWN_SECONDS`
- `ALERT_CHECK_INTERVAL_SECONDS`

## WebSocket Message Shape

Use the documented alert shape unless the user explicitly changes the contract:

```json
{
  "type": "error_spike",
  "count": 25,
  "threshold": 10,
  "window": "5m",
  "timestamp": "2024-01-15T10:23:11Z",
  "message": "25 errors in last 5 minutes (threshold: 10)"
}
```
