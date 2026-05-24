#!/bin/bash
set -e

echo "==> Tạo cấu trúc thư mục log-system..."

# --- Thư mục gốc đã tồn tại (script chạy từ trong đó) ---

# Filebeat
mkdir -p filebeat

# Logstash
mkdir -p logstash/pipeline
mkdir -p logstash/config

# Elasticsearch
mkdir -p elasticsearch/config

# Demo services
mkdir -p services/demo-node
mkdir -p services/demo-go

# Go API server
mkdir -p api-server/handlers
mkdir -p api-server/alerting
mkdir -p api-server/middleware

# Dashboard
mkdir -p dashboard

# Docs
mkdir -p docs

# Logs (volume mount cho demo services)
mkdir -p logs/demo-node
mkdir -p logs/demo-go

echo "==> Tạo file placeholder..."

# .env
cat > .env << 'ENVEOF'
# Elasticsearch
ES_VERSION=8.13.0
ES_PORT=9200
ES_PASSWORD=changeme123

# Logstash
LOGSTASH_PORT=5044

# Go API
API_PORT=8080

# Alerting
ALERT_THRESHOLD=10
ALERT_WINDOW_SECONDS=300
ALERT_COOLDOWN_SECONDS=60
ALERT_CHECK_INTERVAL_SECONDS=10
ENVEOF

# .env.example (commit lên git, không có password thật)
cat > .env.example << 'ENVEOF'
ES_VERSION=8.13.0
ES_PORT=9200
ES_PASSWORD=your_password_here

LOGSTASH_PORT=5044

API_PORT=8080

ALERT_THRESHOLD=10
ALERT_WINDOW_SECONDS=300
ALERT_COOLDOWN_SECONDS=60
ALERT_CHECK_INTERVAL_SECONDS=10
ENVEOF

# .gitignore
cat > .gitignore << 'GITEOF'
.env
logs/
*.log
api-server/bin/
services/demo-go/bin/
GITEOF

# README skeleton
cat > README.md << 'MDEOF'
# Log System — Centralized Logging Platform

## Yêu cầu
- Docker & Docker Compose v2
- 8GB RAM trống (khuyến nghị 12GB)

## Khởi động nhanh
```bash
cp .env.example .env
# Chỉnh sửa .env nếu cần
docker compose up -d
```

## Truy cập
- Dashboard: http://localhost:8080
- Elasticsearch: http://localhost:9200
- API docs: xem docs/api.md

## Kiến trúc
Xem docs/architecture.md
MDEOF

# Filebeat placeholder
touch filebeat/filebeat.yml

# Logstash placeholders
touch logstash/pipeline/logstash.conf
touch logstash/config/logstash.yml

# Elasticsearch placeholder
touch elasticsearch/config/elasticsearch.yml

# Docker Compose placeholder
touch docker-compose.yml

# Go API placeholders
cat > api-server/go.mod << 'GOEOF'
module github.com/yourname/log-system/api-server

go 1.22
GOEOF

touch api-server/main.go
touch api-server/handlers/logs.go
touch api-server/handlers/alerts.go
touch api-server/alerting/engine.go
touch api-server/middleware/cors.go
touch api-server/Dockerfile

# Demo Node.js
cat > services/demo-node/package.json << 'PKGEOF'
{
  "name": "demo-node",
  "version": "1.0.0",
  "description": "Demo service sinh log",
  "main": "index.js",
  "scripts": {
    "start": "node index.js"
  }
}
PKGEOF

touch services/demo-node/index.js
touch services/demo-node/Dockerfile

# Demo Go
cat > services/demo-go/go.mod << 'GOEOF'
module github.com/yourname/log-system/demo-go

go 1.22
GOEOF

touch services/demo-go/main.go
touch services/demo-go/Dockerfile

# Dashboard
touch dashboard/index.html
touch dashboard/app.js
touch dashboard/style.css

# Docs
touch docs/architecture.md
touch docs/api.md
touch docs/deployment.md

echo ""
echo "==> Hoàn thành! Cấu trúc thư mục:"
find . -not -path './.git/*' | sort | sed 's|[^/]*/|  |g;s|  \([^ ]\)|── \1|'
