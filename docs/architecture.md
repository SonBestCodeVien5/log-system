# Kiến trúc hệ thống

## Tổng quan

Hệ thống log tập trung (Centralized Logging Platform) gom log từ nhiều service,
parse và lưu vào Elasticsearch, cung cấp dashboard tìm kiếm và alerting real-time.

## Luồng dữ liệu

```
[Services]
  demo-node (Node.js)  -->  /logs/demo-node/app.log
  demo-go   (Go)       -->  /logs/demo-go/app.log
       |
       v
[Filebeat :5044]
  - Tail log files
  - Buffer khi Logstash down
  - Registry: nhớ đã đọc đến đâu
       |
       v
[Logstash]
  Input:  beats (port 5044)
  Filter: Grok parse text --> JSON có cấu trúc
  Output: Elasticsearch
       |
       v
[Elasticsearch :9200]
  Index: logs-YYYY.MM.DD
  ILM:   giữ 30 ngày
       |
       +-----------> [Go API Server :8080]
       |                  GET  /api/logs
       |                  GET  /api/logs/count
       |                  GET  /api/health
       |                  POST /api/alerts/config
       |                  WS   /ws/alerts
       |
       +-----------> [Alerting Engine - goroutine]
                          Sliding Window 5 phút
                          Deduplication 60s cooldown
                          Dynamic Threshold
                               |
                               v WebSocket
                     [Dashboard :8080]
                          Bảng log + filter
                          Alert banner
```

## Các tầng hệ thống

### Tầng 1 — Demo Services

Sinh log liên tục để test pipeline. Không phải nghiệp vụ thật.

| Service | Công nghệ | Output |
|---|---|---|
| demo-node | Node.js 20 | `/logs/demo-node/app.log` |
| demo-go | Go 1.22 | `/logs/demo-go/app.log` |

Tỉ lệ log: INFO 60% — WARN 25% — ERROR 15%

**Format log bắt buộc:**
```
[2024-01-15T10:23:11Z] [ERROR] [demo-node] Payment gateway timeout after 30s
[2024-01-15T10:23:12Z] [INFO]  [demo-go]   User login successful uid=12345
[2024-01-15T10:23:13Z] [WARN]  [demo-node] Retry attempt 2/3 for order=789
```

### Tầng 2 — Filebeat

Agent nhẹ chạy trên từng server, tail log file và ship tới Logstash.

**Tại sao không ship thẳng vào Logstash từ app?**
- Filebeat có registry — nhớ đã đọc đến dòng nào, không mất log khi Logstash tạm down
- Footprint nhỏ (~50MB) so với Logstash (~500MB), không nặng server

### Tầng 3 — Logstash

Parse log text thô thành JSON có cấu trúc bằng Grok pattern.

```
Input raw:   "[2024-01-15T10:23:11Z] [ERROR] [demo-node] Payment failed"
Output JSON: {
  "@timestamp": "2024-01-15T10:23:11Z",
  "level":      "ERROR",
  "service":    "demo-node",
  "message":    "Payment failed"
}
```

### Tầng 4 — Elasticsearch

Lưu trữ và search log. Index theo ngày để tối ưu query và lifecycle.

**Tại sao không dùng MySQL/PostgreSQL?**

| Tiêu chí | Elasticsearch | MySQL |
|---|---|---|
| Index type | Inverted index | B-tree |
| Full-text search | Tự nhiên, rất nhanh | Chậm, cần LIKE |
| Time-range query | `now-5m` built-in | Tính timestamp thủ công |
| Scale | Horizontal dễ | Vertical tốn kém |

### Tầng 5 — Go API Server

REST API + WebSocket viết bằng Go + gin, trung gian giữa ES và Dashboard.

**Tại sao Go thay vì Spring Boot Java (yêu cầu gốc):**
- Docker, Kubernetes, Prometheus, Filebeat đều viết bằng Go — đúng domain
- Goroutine native phù hợp với alerting engine concurrent
- RAM 20–50MB so với Spring 300–500MB
- Code tường minh, không có annotation magic

## Alerting Engine

### Sliding Window
Quét mỗi 10 giây, nhìn lùi 5 phút — không bỏ sót spike lỗi ngắn.

```
t=0s:  đếm ERROR trong [t-300s, t] = 3  --> bình thường
t=10s: đếm ERROR trong [t-300s, t] = 15 --> ALERT!
t=20s: đếm ERROR trong [t-300s, t] = 14 --> dedup, không bắn lại
```

### Alert Deduplication
Không gửi alert trùng trong cooldown period — tránh Alert Fatigue.

### Dynamic Threshold
Dashboard gửi ngưỡng mới về API, goroutine alerting cập nhật ngay không cần restart.
Dùng `sync.RWMutex` để đảm bảo an toàn khi 2 goroutine đọc/ghi đồng thời.

## Cấu trúc thư mục

```
log-system/
├── AGENTS.md
├── docker-compose.yml
├── .env / .env.example
├── filebeat/filebeat.yml
├── logstash/
│   ├── pipeline/logstash.conf
│   └── config/logstash.yml
├── elasticsearch/config/elasticsearch.yml
├── services/
│   ├── demo-node/index.js
│   └── demo-go/main.go
├── api-server/
│   ├── main.go
│   ├── handlers/logs.go
│   ├── handlers/alerts.go
│   ├── alerting/engine.go
│   └── middleware/cors.go
├── dashboard/
│   ├── index.html
│   ├── app.js
│   └── style.css
├── logs/
│   ├── demo-node/
│   └── demo-go/
└── docs/
```
