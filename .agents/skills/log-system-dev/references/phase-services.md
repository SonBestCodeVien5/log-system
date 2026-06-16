# Phase: Demo Services

Use this context for `services/demo-node` and `services/demo-go`.

## Purpose

Demo services exist only to generate logs for testing the pipeline.

## Required Log Format

Do not change this format:

```json
{"timestamp":"2024-01-15T10:23:11Z","level":"ERROR","service":"demo-node","message":"Payment gateway timeout after 30s","metadata":{"order_id":"789"}}
```

Rules:

- Each line is one complete JSON object.
- `timestamp` is ISO8601 UTC.
- `level` is uppercase: `INFO`, `WARN`, or `ERROR`.
- `service` has no spaces.
- `message` is free text for dashboard display.
- `metadata` is an object for extra fields.

## Runtime Behavior

- Write Node logs to `/logs/demo-node/app.log`.
- Write Go logs to `/logs/demo-go/app.log`.
- Generate INFO around 60%, WARN around 25%, ERROR around 15%.
- Emit one log every 1-3 seconds by default.
- Include a way to increase ERROR frequency for alert testing if the user asks.

## Module Paths

- Go demo module: `github.com/SonBestCodeVien5/log-system/demo-go`.
