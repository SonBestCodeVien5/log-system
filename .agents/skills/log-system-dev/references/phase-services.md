# Phase: Demo Services

Use this context for `services/demo-node` and `services/demo-go`.

## Purpose

Demo services exist only to generate logs for testing the pipeline.

## Required Log Format

Do not change this format:

```text
[2024-01-15T10:23:11Z] [ERROR] [demo-node] Payment gateway timeout after 30s
```

Rules:

- Timestamp is ISO8601 UTC.
- Level is uppercase: `INFO`, `WARN`, or `ERROR`.
- Service name has no spaces.
- Message is free text after the third bracket.

## Runtime Behavior

- Write Node logs to `/logs/demo-node/app.log`.
- Write Go logs to `/logs/demo-go/app.log`.
- Generate INFO around 60%, WARN around 25%, ERROR around 15%.
- Emit one log every 1-3 seconds by default.
- Include a way to increase ERROR frequency for alert testing if the user asks.

## Module Paths

- Go demo module: `github.com/SonBestCodeVien5/log-system/demo-go`.
